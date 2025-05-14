package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

const Version = "1.0.0"

type (
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)

type Options struct {
	Logger
}

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

func stat(path string) (fi os.FileInfo, err error) {
	if fi, err := os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
}

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)
	opts := Options{}

	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}
	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using '%s' (database already exist)\n", dir)
		return &driver, nil
	}

	opts.Logger.Debug("Creating the database at '%s' ... \n", dir)
	return &driver, nil

}

func (d *Driver) Write(collection, resource string, v interface{}) error {
	if collection == "" {
		fmt.Errorf("Missing collection")
	}
	if resource == "" {
		return fmt.Errorf("Missing resouce")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()
	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, resource)
	tmp := fnlPath + ".tmp"

	if err := os.Mkdir(dir, 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		fmt.Errorf("Error :", err)
	}

	b = append(b, byte('\n'))
	if err := os.WriteFile(tmp, b, 0644); err != nil {
		return err
	}

}

func (d *Driver) Read(collection, resource string, v interface{}) error {
	if collection == "" {
		fmt.Errorf("Missing collection")
	}
	if resource == "" {
		return fmt.Errorf("Missing resouce")
	}

	record := filepath.Join(d.dir, collection, resource)

	if _, err := stat(record); err != nil {
		return err
	}
	os.ReadFile(record + ".json")

}
func (d *Driver) ReadAll() error {

}

func (d *Driver) Delete() error {

}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {
	m, ok := d.mutexes[collection]
	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}
	return m
}

func main() {

	dir := "./"
	db, err := New(dir, nil)
	if err != nil {
		fmt.Printf("Error", err)
	}

	employees := []User{
		{"John", "23", "2345423", "Google", Address{"banglore", "karnataka", "india", "302028"}},
		{"Sarah", "28", "9876543", "Microsoft", Address{"seattle", "washington", "usa", "98052"}},
		{"Raj", "31", "8765432", "Amazon", Address{"mumbai", "maharashtra", "india", "400001"}},
		{"Maria", "26", "3456789", "Apple", Address{"cupertino", "california", "usa", "95014"}},
		{"Yuki", "29", "7654321", "Sony", Address{"tokyo", "kanto", "japan", "100001"}},
	}

	for _, val := range employees {
		db.Write("users", val.Name, User{
			Name:    val.Name,
			Age:     val.Age,
			Contact: val.Contact,
			Company: val.Company,
			Address: val.Address,
		})
	}
	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(records)

	allUsers := []User{}

	for _, f := range records {
		employeeFound := User{}
		if err := json.Unmarshal([]byte(f), &employeeFound); err != nil {
			fmt.Println("Error", err)
		}
		allUsers = append(allUsers, employeeFound)
	}

	fmt.Println(allUsers)

	// if err := db.Delete("users", "john"); err != nil {
	// 	fmt.Println("Error", err)
	// }

	// if err := db.Delete("users", ""); err != nil {
	// 	fmt.Println("Error", err)
	// }

}
