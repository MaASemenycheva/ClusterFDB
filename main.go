package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
	"unsafe"

	"fdb/models"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	db      fdb.Database
	jsonKey = append([]byte{255, 255}, []byte("/status/json")...)
)

func stringInsert1(db fdb.Database) {
	time.Sleep(1 * time.Millisecond)
	_, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, e error) {
		tr.Set(fdb.Key("str "+strconv.Itoa(0)), []byte("str "+strconv.Itoa(0)))
		return
	})
	if err != nil {
		log.Fatalf("Unable to set FDB database value (%v)", err)
	}
}

func stringSelect1(db fdb.Database) {
	time.Sleep(1 * time.Millisecond)
	//ret
	_, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, e error) {
		ret = tr.Get(fdb.Key("str " + strconv.Itoa(0))).MustGet()
		return
	})
	if err != nil {
		log.Fatalf("Unable to read FDB database value (%v)", err)
	}
	//v := ret.([]byte)
	//fmt.Printf("string row = , %s\n", string(v))
}

func stringInsert10000(db fdb.Database) {
	for i := 0; i < 10000; i++ {
		time.Sleep(1 * time.Millisecond)
		_, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, e error) {
			tr.Set(fdb.Key("str "+strconv.Itoa(i)), []byte("str "+strconv.Itoa(i)))
			return
		})
		if err != nil {
			log.Fatalf("Unable to set FDB database value (%v)", err)
		}
	}
}

func stringSelect10000(db fdb.Database) {
	for i := 0; i < 10000; i++ {
		time.Sleep(1 * time.Millisecond)
		//ret
		ret, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, e error) {
			ret = tr.Get(fdb.Key("str " + strconv.Itoa(i))).MustGet()
			return
		})
		if err != nil {
			log.Fatalf("Unable to read FDB database value (%v)", err)
		}
		v := ret.([]byte)
		fmt.Printf("string row = , %s\n", string(v))
	}
}

func IntToByteArray(num int64) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}

func ByteArrayToInt(arr []byte) int64 {
	val := int64(0)
	size := len(arr)
	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}

func intInsert1(db fdb.Database) {
	time.Sleep(1 * time.Millisecond)
	_, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, e error) {
		num := int64(7)
		//fmt.Println("Original number:", num)
		// integer to byte array
		byteArr := IntToByteArray(num)
		//fmt.Println("Byte Array", byteArr)
		tr.Set(fdb.Key("uint"+strconv.Itoa(0)), byteArr)
		return
	})
	if err != nil {
		log.Fatalf("Unable to set FDB database value (%v)", err)
	}
}

func intSelect1(db fdb.Database) {
	time.Sleep(1 * time.Millisecond)
	//ret
	_, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, e error) {
		ret = tr.Get(fdb.Key("uint" + strconv.Itoa(0))).MustGet()
		return
	})
	if err != nil {
		log.Fatalf("Unable to read FDB database value (%v)", err)
	}
	//v := ret.([]byte)
	//fmt.Printf("int_number = , %s\n", BytesToInt(v))
	// byte array to integer again
	//numAgain := ByteArrayToInt(v)
	//fmt.Println("Converted number:", numAgain)
}

func intInsert10000(db fdb.Database) {
	for i := 0; i < 10000; i++ {
		time.Sleep(1 * time.Millisecond)
		_, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, e error) {
			num := int64(7)
			//fmt.Println("Original number:", num)
			// integer to byte array
			byteArr := IntToByteArray(num)
			//fmt.Println("Byte Array", byteArr)
			tr.Set(fdb.Key("uint"+strconv.Itoa(i)), byteArr)
			return
		})
		if err != nil {
			log.Fatalf("Unable to set FDB database value (%v)", err)
		}
	}
}

func intSelect10000(db fdb.Database) {
	for i := 0; i < 10000; i++ {
		time.Sleep(1 * time.Millisecond)
		//ret
		_, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, e error) {
			ret = tr.Get(fdb.Key("uint" + strconv.Itoa(i))).MustGet()
			return
		})
		if err != nil {
			log.Fatalf("Unable to read FDB database value (%v)", err)
		}
		//v := ret.([]byte)
		//fmt.Printf("int number = , %s\n", BytesToInt(v))
		// byte array to integer again
		//numAgain := ByteArrayToInt(v)
		//fmt.Println("Converted number:", numAgain)
	}
}

func Float32ToByte(f float32) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, f)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func Float32frombytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func floatInsert1(db fdb.Database) {
	time.Sleep(1 * time.Millisecond)
	_, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, e error) {
		bytes := Float32ToByte(math.Pi)
		tr.Set(fdb.Key("float"+strconv.Itoa(0)), bytes)
		return
	})
	if err != nil {
		log.Fatalf("Unable to set FDB database value (%v)", err)
	}
}

func floatSelect1(db fdb.Database) {
	time.Sleep(1 * time.Millisecond)
	//ret
	_, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, e error) {
		ret = tr.Get(fdb.Key("float" + strconv.Itoa(0))).MustGet()
		return
	})
	if err != nil {
		log.Fatalf("Unable to read FDB database value (%v)", err)
	}
	//v := ret.([]byte)
	//float := Float32frombytes(v)
	//fmt.Println(float)
}

func floatInsert10000(db fdb.Database) {
	for i := 0; i < 10000; i++ {
		time.Sleep(1 * time.Millisecond)
		_, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, e error) {
			bytes := Float32ToByte(math.Pi + 0.1*float32(i))
			tr.Set(fdb.Key("float"+strconv.Itoa(i)), bytes)
			return
		})
		if err != nil {
			log.Fatalf("Unable to set FDB database value (%v)", err)
		}
	}
}

func floatSelect10000(db fdb.Database) {
	for i := 0; i < 10000; i++ {
		time.Sleep(1 * time.Millisecond)
		//ret
		_, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, e error) {
			ret = tr.Get(fdb.Key("float" + strconv.Itoa(i))).MustGet()
			return
		})
		if err != nil {
			log.Fatalf("Unable to read FDB database value (%v)", err)
		}
		//v := ret.([]byte)
		//float := Float32frombytes(v)
		//fmt.Println(float)
	}
}

func vectorInsert1(db fdb.Database) {
	time.Sleep(1 * time.Millisecond)
	_, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, e error) {
		var args = `[1.1, 7.5, 1.7]`
		tr.Set(fdb.Key("vector"+strconv.Itoa(0)), []byte(args))
		return
	})
	if err != nil {
		log.Fatalf("Unable to set FDB database value (%v)", err)
	}
}

func vectorSelect1(db fdb.Database) {
	time.Sleep(1 * time.Millisecond)
	ret, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, e error) {
		ret = tr.Get(fdb.Key("vector" + strconv.Itoa(0))).MustGet()
		return
	})
	if err != nil {
		log.Fatalf("Unable to read FDB database value (%v)", err)
	}
	v := ret.([]byte)
	//fmt.Printf("vector, %s\n", string(v))
	var x []interface{}
	err1 := json.Unmarshal([]byte(v), &x)
	if err1 != nil {
		log.Fatalf("%s", err1.Error())
		panic(fmt.Sprintf("%s", err1.Error()))
	}
	//for _, arg := range x {
	//    t := reflect.TypeOf(arg).Kind().String()
	//    v := reflect.ValueOf(arg)
	//    if t == "int64" {
	//        fmt.Printf("int64 %v\n", v.Int())
	//    }
	//    if t == "float64" {
	//        fmt.Printf("float64 %v\n", v.Float())
	//    }
	//    if t == "string" {
	//        fmt.Printf("string %v\n", v.String())
	//    }
	//    if t == "bool" {
	//        fmt.Printf("bool %v\n", v.Bool())
	//    }
	//}
}

func vectorInsert10000(db fdb.Database) {
	for i := 0; i < 10000; i++ {
		time.Sleep(1 * time.Millisecond)
		_, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, e error) {
			var args = `[1.1, 7.5, 1.7]`
			tr.Set(fdb.Key("vector"+strconv.Itoa(i)), []byte(args))
			return
		})
		if err != nil {
			log.Fatalf("Unable to set FDB database value (%v)", err)
		}
	}
}

func vectorSelect10000(db fdb.Database) {
	for i := 0; i < 10000; i++ {
		time.Sleep(1 * time.Millisecond)
		_, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, e error) {
			ret = tr.Get(fdb.Key("vector" + strconv.Itoa(i))).MustGet()
			return
		})
		if err != nil {
			log.Fatalf("Unable to read FDB database value (%v)", err)
		}
		//v := ret.([]byte)
		//fmt.Printf("vector, %s\n", string(v))
		//var x []interface{}
		//err1 := json.Unmarshal([]byte(v), &x)
		//if err1 != nil {
		//	log.Fatalf("%s", err1.Error())
		//	panic(fmt.Sprintf("%s", err1.Error()))
		//}
		//for _, arg := range x {
		//    t := reflect.TypeOf(arg).Kind().String()
		//    v := reflect.ValueOf(arg)
		//    if t == "int64" {
		//        fmt.Printf("int64 %v\n", v.Int())
		//    }
		//    if t == "float64" {
		//        fmt.Printf("float64 %v\n", v.Float())
		//    }
		//    if t == "string" {
		//        fmt.Printf("string %v\n", v.String())
		//    }
		//    if t == "bool" {
		//        fmt.Printf("bool %v\n", v.Bool())
		//    }
		//}
	}
}

func main() {

	apiVersion, err := strconv.Atoi(getEnv("FDB_API_VERSION", "620"))
	if err != nil {
		log.Fatal("cannot parse FDB_API_VERSION from env")
	}
	// Different API versions may expose different runtime behaviors.
	fdb.MustAPIVersion(apiVersion)

	clusterFile := getEnv("FDB_CLUSTER_FILE", "/var/fdb/data/fdb.cluster")

	if _, exists := os.LookupEnv("FDB_CREATE_CLUSTER_FILE"); exists {
		createClusterFile()
	}

	fmt.Println("opening cluster file at", clusterFile)
	data, err := ioutil.ReadFile(clusterFile)
	if err != nil {
		log.Fatalf("cannot read cluster file: %+v", err)
	}
	fmt.Println(string(data))

	// Open the default database from the system cluster
	db = fdb.MustOpenDatabase(clusterFile)

	exportWorkload, err := strconv.ParseBool(getEnv("FDB_EXPORT_WORKLOAD", "true"))
	if err != nil {
		log.Fatal("cannot parse FDB_EXPORT_WORLOAD from env")
	}
	exportDatabaseStatus, err := strconv.ParseBool(getEnv("FDB_EXPORT_DATABASE_STATUS", "true"))
	if err != nil {
		log.Fatal("cannot parse FDB_EXPORT_DATABASE_STATUS from env")
	}

	exportConfiguration, err := strconv.ParseBool(getEnv("FDB_EXPORT_CONFIGURATION", "true"))
	if err != nil {
		log.Fatal("cannot parse FDB_EXPORT_CONFIGURATION from env")
	}
	exportProcesses, err := strconv.ParseBool(getEnv("FDB_EXPORT_PROCESSES", "true"))
	if err != nil {
		log.Fatal("cannot parse FDB_EXPORT_PROCESSES from env")
	}

	listenTo := getEnv("FDB_METRICS_LISTEN", ":8080")
	refreshEvery, err := strconv.Atoi(getEnv("FDB_METRICS_EVERY", "10"))
	if err != nil {
		log.Fatal("cannot parse FDB_METRICS_EVERY from env")
	}

	ticker := time.NewTicker(time.Duration(refreshEvery) * time.Second)
	go func() {
		for range ticker.C {
			//Call the periodic function here.
			models, err := retrieveMetrics()
			if err != nil {
				fmt.Errorf("cannot retrieve metrics from FDB: (%v)", err)
				continue
			}

			fmt.Println("µs - Microsecond (1µs=0.001ms)")
			fmt.Println("ns - Nanosecond (1ns=0.000001ms)")
			fmt.Println("ms - Millisecond (1ms=1ms)")

			startStringInsert1 := time.Now()
			go stringInsert1(db)
			durationStringInsert1 := time.Since(startStringInsert1)
			fmt.Println("stringInsert1 time = ", durationStringInsert1, "\n")

			startStringSelect1 := time.Now()
			go stringSelect1(db)
			durationStringSelect1 := time.Since(startStringSelect1)
			fmt.Println("stringSelect1 time = ", durationStringSelect1, "\n")

			startStringInsert10000 := time.Now()
			go stringInsert10000(db)
			durationStringInsert10000 := time.Since(startStringInsert10000)
			fmt.Println("stringStringInsert10000 time = ", durationStringInsert10000, "\n")

			startStringSelect10000 := time.Now()
			go stringSelect10000(db)
			durationStringSelect10000 := time.Since(startStringSelect10000)
			fmt.Println("stringSelect10000 time = ", durationStringSelect10000, "\n")

			startIntInsert1 := time.Now()
			go intInsert1(db)
			durationIntInsert1 := time.Since(startIntInsert1)
			fmt.Println("intInsert1 time = ", durationIntInsert1, "\n")

			startIntSelect1 := time.Now()
			go intSelect1(db)
			durationIntSelect1 := time.Since(startIntSelect1)
			fmt.Println("intSelect1 time = ", durationIntSelect1, "\n")

			startIntInsert10000 := time.Now()
			go intInsert10000(db)
			durationIntInsert10000 := time.Since(startIntInsert10000)
			fmt.Println("intInsert10000 time = ", durationIntInsert10000, "\n")

			startIntSelect10000 := time.Now()
			go intSelect10000(db)
			durationIntSelect10000 := time.Since(startIntSelect10000)
			fmt.Println("intSelect10000 time = ", durationIntSelect10000, "\n")

			startFloatInsert1 := time.Now()
			go floatInsert1(db)
			durationFloatInsert1 := time.Since(startFloatInsert1)
			fmt.Println("floatInsert1 time = ", durationFloatInsert1, "\n")

			startFloatSelect1 := time.Now()
			go floatSelect1(db)
			durationFloatSelect1 := time.Since(startFloatSelect1)
			fmt.Println("floatSelect1 time = ", durationFloatSelect1, "\n")

			startFloatInsert10000 := time.Now()
			go floatInsert10000(db)
			durationFloatInsert10000 := time.Since(startFloatInsert10000)
			fmt.Println("floatInsert10000 time = ", durationFloatInsert10000, "\n")

			startFloatSelect10000 := time.Now()
			go floatSelect10000(db)
			durationFloatSelect10000 := time.Since(startFloatSelect10000)
			fmt.Println("floatSelect10000 time = ", durationFloatSelect10000, "\n")

			startVectorInsert1 := time.Now()
			go vectorInsert1(db)
			durationVectorInsert1 := time.Since(startVectorInsert1)
			fmt.Println("vectorInsert1 time = ", durationVectorInsert1, "\n")

			startVectorSelect1 := time.Now()
			go vectorSelect1(db)
			durationVectorSelect1 := time.Since(startVectorSelect1)
			fmt.Println("vectorSelect1 time = ", durationVectorSelect1, "\n")

			startVectorInsert10000 := time.Now()
			go vectorInsert10000(db)
			durationVectorInsert10000 := time.Since(startVectorInsert10000)
			fmt.Println("vectorInsert10000 time = ", durationVectorInsert10000, "\n")

			startVectorSelect10000 := time.Now()
			go vectorSelect10000(db)
			durationVectorSelect10000 := time.Since(startVectorSelect10000)
			fmt.Println("vectorSelect10000 time = ", durationVectorSelect10000, "\n")

			fmt.Println("retrieved data")

			if exportWorkload {
				models.ExportWorkload()
			}

			if exportDatabaseStatus {
				models.ExportDatabaseStatus()
			}

			if exportConfiguration {
				models.ExportConfiguration()
			}

			if exportProcesses {
				models.ExportProcesses()
			}
		}
	}()

	r := prometheus.NewRegistry()
	models.Register(r) // Role model

	// [...] update metrics within a goroutine
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)

	log.Fatal(http.ListenAndServe(listenTo, nil))
}

func retrieveMetrics() (*models.FDBStatus, error) {

	jsonRaw, err := db.ReadTransact(func(tr fdb.ReadTransaction) (interface{}, error) {
		return tr.Get(fdb.Key(jsonKey)).Get()
	})

	if err != nil {
		return nil, errors.Wrap(err, "cannot get status")
	}

	var status models.FDBStatus
	err = json.Unmarshal([]byte(jsonRaw.([]byte)), &status)
	if err != nil {
		return nil, errors.Wrap(err, "cannot decode json")
	}
	return &status, nil
}

// getEnv is wrapping os.getenv with a fallback
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func createClusterFile() {
	cmd := exec.Command("/create_cluster_file.bash")

	fmt.Printf("Running command 'create_cluster_file' and waiting for it to finish...\n")
	stdoutStderr, err := cmd.CombinedOutput()
	fmt.Printf("%s\n", stdoutStderr)
	if err != nil {
		log.Fatal(err)
	}
}
