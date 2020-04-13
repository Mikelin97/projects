// package proj2

// // You MUST NOT change what you import.  If you add ANY additional
// // imports it will break the autograder, and we will be Very Upset.

// import (
// 	"testing"
// 	"reflect"
// 	"github.com/cs161-staff/userlib"
// 	_ "encoding/json"
// 	_ "encoding/hex"
// 	_ "github.com/google/uuid"
// 	_ "strings"
// 	_ "errors"
// 	_ "strconv"
// )

// func clear() {
// 	// Wipes the storage so one test does not affect another
// 	userlib.DatastoreClear()
// 	userlib.KeystoreClear()
// }

// func TestInit(t *testing.T) {
// 	clear()
// 	t.Log("Initialization test")

// 	// You can set this to false!
// 	userlib.SetDebugStatus(true)

// 	u, err := InitUser("alice", "fubar")
// 	if err != nil {
// 		// t.Error says the test fails
// 		t.Error("Failed to initialize user", err)
// 		return
// 	}
// 	// t.Log() only produces output if you run with "go test -v"
// 	t.Log("Got user", &u)
// 	// If you want to comment the line above,
// 	// write _ = u here to make the compiler happy
// 	// You probably want many more tests here.
// }

// func TestGet(t *testing.T) {
// 	// Test the Get Function
// 	clear()
// 	t.Log("Initialization test")

// 	// You can set this to false!
// 	userlib.SetDebugStatus(true)

// 	u, err := InitUser("mike", "fubar")
// 	if err != nil {
// 		// t.Error says the test fails
// 		t.Error("Failed to initialize user", err)
// 		return
// 	}

// 	// t.Log() only produces output if you run with "go test -v"
// 	t.Log("Got user", &u)
// 	// If you want to comment the line above,
// 	// write _ = u here to make the compiler happy
// 	// You probably want many more tests here.
// 	u1, err := GetUser("mike", "fubar")
// 	t.Log("Check GetUser", &u1)
// 	// t.Log("Username?", u1.Username)

// }

// func TestGetEmptyUser(t *testing.T) {
// 	// Test the Get Function
// 	clear()
// 	t.Log("Initialization test")

// 	// You can set this to false!
// 	userlib.SetDebugStatus(true)

// 	// t.Log() only produces output if you run with "go test -v"
// 	// If you want to comment the line above,
// 	// write _ = u here to make the compiler happy
// 	// You probably want many more tests here.
// 	u1, _ := GetUser("mike", "fubar")
	
// 	t.Log("Check GetUser", &u1)
// 	// t.Log("Username?", u1.Username)

// }

// func TestGetSalt(t *testing.T) {
// 	// Test the Get Function as well 
// 	clear()
// 	t.Log("Initialization test")

// 	// You can set this to false!
// 	userlib.SetDebugStatus(true)

// 	u, err := InitUser("mike", "fubar")
// 	if err != nil {
// 		// t.Error says the test fails
// 		t.Error("Failed to initialize user", err)
// 		return
// 	}

// 	// t.Log() only produces output if you run with "go test -v"
// 	t.Log("Got user", &u)
// 	// t.Log("Test1", u.Test1)
// 	// t.Log("Test2", u.Test2)
// 	// If you want to comment the line above,
// 	// write _ = u here to make the compiler happy
// 	// You probably want many more tests here.
// 	u1, err := GetUser("mike", "fubar")
// 	t.Log("Check GetUser", &u1)
// 	t.Log("PasswordKey?", u1.PasswordKey)

// }


// func TestStorageOnly(t *testing.T) {
// 	clear()
// 	u, err := InitUser("alice", "fubar")
// 	if err != nil {
// 		t.Error("Failed to initialize user", err)
// 		return
// 	}

// 	v := []byte("This is a test")
// 	u.StoreFile("file1", v)

// 	v2, err2 := u.LoadFile("file1")
// 	if err2 != nil {
// 		t.Error("Failed to upload and download", err2)
// 		return
// 	}
// 	if !reflect.DeepEqual(v, v2) {
// 		t.Error("Downloaded file is not the same", v, v2)
// 		return
// 	}
// }

// func TestStorage(t *testing.T) {
// 	clear()
// 	u, err := InitUser("alice", "fubar")
// 	if err != nil {
// 		t.Error("Failed to initialize user", err)
// 		return
// 	}

// 	v := []byte("This is a test")
// 	u.StoreFile("file1", v)

// 	v2, err2 := u.LoadFile("file1")
// 	if err2 != nil {
// 		t.Error("Failed to upload and download", v2)
// 		return
// 	}
// 	if !reflect.DeepEqual(v, v2) {
// 		t.Error("Downloaded file is not the same", v, v2)
// 		return
// 	}
// }

// // func TestAppend(t *testing.T) {
// // 	clear()
// // 	u, err := InitUser("alice", "fubar")
// // 	if err != nil {
// // 		t.Error("Failed to initialize user", err)
// // 		return
// // 	}

// // 	v := []byte("This is a test")
// // 	u.StoreFile("file1", v)
// // 	v1 := []byte("another test")
// // 	u.AppendFile("file1", v1)
// // 	v3 := append(v, v1...)

// // 	v2, err2 := u.LoadFile("file1")
// // 	if err2 != nil {
// // 		t.Error("Failed to upload and download", v2)
// // 		return
// // 	}
// // 	if !reflect.DeepEqual(v2, v3) {
// // 		t.Error("Downloaded file is not the same", v2, v3)
// // 		return
// // 	}
// // }

// func TestUser_AppendFile(t *testing.T) {
// 	u, err := GetUser("alice", "fubar")
// 	v,_ := u.LoadFile("file1")
// 	if err != nil {
// 		t.Error("Failed to download the file from alice", err)
// 		return
// 	}
// 	err = u.AppendFile("file1", []byte("helloworld"))
// 	if err != nil {
// 		t.Error("Failed to append file", err)
// 		return
// 	}
// 	v2,_ := u.LoadFile("file1")
// 	if err != nil {
// 		t.Error("Failed to download the file from alice", err)
// 		return
// 	}
	
// 	if !reflect.DeepEqual(append(v, []byte("helloworld")...), v2) {
// 		t.Error("Appending wrong", v, v2)
// 		return
// 	}
// }

// func TestUser_AppendMultipleFile(t *testing.T) {
// 	u, err := GetUser("alice", "fubar")
// 	v,_ := u.LoadFile("file1")
// 	if err != nil {
// 		t.Error("Failed to download the file from alice", err)
// 		return
// 	}
// 	err = u.AppendFile("file1", []byte("helloworld"))
// 	if err != nil {
// 		t.Error("Failed to append file", err)
// 		return
// 	}

// 	err = u.AppendFile("file1", []byte("helloworld"))
// 	if err != nil {
// 		t.Error("Failed to append file", err)
// 		return
// 	}

// 	v2,_ := u.LoadFile("file1")
// 	if err != nil {
// 		t.Error("Failed to download the file from alice", err)
// 		return
// 	}
// 	v3 := append(v, []byte("helloworld")...)
// 	if !reflect.DeepEqual(append(v3, []byte("helloworld")...), v2) {
// 		t.Error("Appending wrong", v3, v2)
// 		return
// 	}
// }

// func TestInvalidFile(t *testing.T) {
// 	clear()
// 	u, err := InitUser("alice", "fubar")
// 	if err != nil {
// 		t.Error("Failed to initialize user", err)
// 		return
// 	}

// 	_, err2 := u.LoadFile("this file does not exist")
// 	if err2 == nil {
// 		t.Error("Downloaded a nonexistent file", err2)
// 		return
// 	}
// }


// func TestShare(t *testing.T) {
// 	clear()
// 	u, err := InitUser("alice", "fubar")
// 	if err != nil {
// 		t.Error("Failed to initialize user", err)
// 		return
// 	}
// 	u2, err2 := InitUser("bob", "foobar")
// 	if err2 != nil {
// 		t.Error("Failed to initialize bob", err2)
// 		return
// 	}

// 	v := []byte("This is a test")
// 	u.StoreFile("file1", v)
	

// 	v, err = u.LoadFile("file1")
// 	if err != nil {
// 		t.Error("Failed to download the file from alice", err)
// 		return
// 	}

// 	magic_string, err := u.ShareFile("file1", "bob")
// 	if err != nil {
// 		t.Error("Failed to share the a file", err)
// 		return
// 	}
// 	err = u2.ReceiveFile("file2", "alice", magic_string)
// 	if err != nil {
// 		t.Error("Failed to receive the share message", err)
// 		return
// 	}

// 	v2, err := u2.LoadFile("file2")
// 	if err != nil {
// 		t.Error("Failed to download the file after sharing", err)
// 		return
// 	}
// 	if !reflect.DeepEqual(v, v2) {
// 		t.Error("Shared file is not the same", v, v2)
// 		return
// 	}

// }

// func TestShareandAppend(t *testing.T) {
// 	clear()
// 	userlib.SetDebugStatus(true)
// 	u, err := InitUser("alice", "fubar")
// 	if err != nil {
// 		t.Error("Failed to initialize user", err)
// 		return
// 	}
// 	u2, err2 := InitUser("bob", "foobar")
// 	if err2 != nil {
// 		t.Error("Failed to initialize bob", err2)
// 		return
// 	}

// 	v := []byte("This is a test")
// 	u.StoreFile("file1", v)
	
	

// 	v, err = u.LoadFile("file1")
// 	if err != nil {
// 		t.Error("Failed to download the file from alice", err)
// 		return
// 	}

// 	magic_string, err := u.ShareFile("file1", "bob")
// 	if err != nil {
// 		t.Error("Failed to share the a file", err)
// 		return
// 	}
// 	err = u2.ReceiveFile("file2", "alice", magic_string)
// 	if err != nil {
// 		t.Error("Failed to receive the share message", err)
// 		return
// 	}

// 	v2, err := u2.LoadFile("file2")
// 	if err != nil {
// 		t.Error("Before Bob append Failed to download the file after sharing", err)
// 		return
// 	}
// 	if !reflect.DeepEqual(v, v2) {
// 		t.Error("Shared file is not the same", v, v2)
// 		return
// 	}

// 	err = u2.AppendFile("file2", []byte("some random testing"))
// 	if err != nil {
// 		t.Error("BOB Failed to append file", err)
// 		return
// 	}
// 	v3 := append(v, []byte("some random testing")...)
// 	v4, err := u.LoadFile("file1")
// 	if err != nil {
// 		t.Error(" After Bob append, ALice Failed to download the file after sharing", err)
// 		return
// 	}
// 	if !reflect.DeepEqual(v3, v4) {
// 		t.Error("Alice's file is not updated after BOB's append", v3, v4)
// 		return
// 	}

// }



// func TestMultipleInstance(t *testing.T) {
// 	clear()


// 	u, err := InitUser("mike", "fubar")
// 	if err != nil {
// 		// t.Error says the test fails
// 		t.Error("Failed to initialize user", err)
// 		return
// 	}

// 	// t.Log() only produces output if you run with "go test -v"
// 	t.Log("Got user", &u)
// 	// If you want to comment the line above,
// 	// write _ = u here to make the compiler happy
// 	// You probably want many more tests here.
// 	u1, err := GetUser("mike", "fubar")
// 	t.Log("Get User instance 1", &u1)

// 	v := []byte("This is a test")
// 	u1.StoreFile("file1", v)


// 	u2, err := GetUser("mike", "fubar")
// 	t.Log("Get User instance 2", &u2)

// 	v2, err2 := u2.LoadFile("file1")
// 	if err2 != nil {
// 		t.Error("Failed to upload and download", v2)
// 		return
// 	}
// 	if !reflect.DeepEqual(v, v2) {
// 		t.Error("Downloaded file is not the same", v, v2)
// 		return
// 	}

// 	err = u2.AppendFile("file1", []byte("helloworld"))
// 	if err != nil {
// 		t.Error("Failed to append file", err)
// 		return
// 	}
// 	v3,_ := u1.LoadFile("file1")
// 	if err != nil {
// 		t.Error("Failed to download the file from alice", err)
// 		return
// 	}
	
// 	if !reflect.DeepEqual(append(v, []byte("helloworld")...), v3) {
// 		t.Error("Appending wrong", v, v2)
// 		return
// 	}

	
	
// }

// // func TestUser_RevokeFile(t *testing.T) {
// // 	u, err := GetUser("alice", "fubar")
// // 	if err != nil {
// // 		t.Error("Failed to reload user", err)
// // 		return
// // 	}

// // 	var v, v2, v3 []byte

// // 	v, err = u.LoadFile("file1")
// // 	if err != nil {
// // 		t.Error("Failed to download the file from alice", err)
// // 		return
// // 	}

// // 	u2, err := GetUser("bob", "foobar")
// // 	if err != nil {
// // 		t.Error("Failed to reload user", err)
// // 		return
// // 	}

// // 	v2, err = u2.LoadFile("file2")
// // 	if err != nil {
// // 		t.Error("Failed to download the file from bob", err)
// // 		return
// // 	}
// // 	if !reflect.DeepEqual(v, v2) {
// // 		t.Error("Shared file is initially not the same", v, v2)
// // 		return
// // 	}

// // 	err = u.RevokeFile("file1", "bob")
// // 	if err != nil {
// // 		t.Error("Failed to revoke file", err)
// // 		return
// // 	}

// // 	v3, err = u2.LoadFile("file2")
	
// // 	if err == nil {
// // 		t.Error("Can still download the file from bob", err)
// // 		return
// // 	}
// // 	if reflect.DeepEqual(v, v3) {
// // 		t.Error("Shared file is still the same", v, v3)
// // 		return
// // 	}


// // }


package proj2

// You MUST NOT change what you import.  If you add ANY additional
// imports it will break the autograder, and we will be Very Upset.

import (
	"testing"
	"reflect"
	"github.com/cs161-staff/userlib"
	_ "encoding/json"
	_ "encoding/hex"
	_ "github.com/google/uuid"
	_ "strings"
	_ "errors"
	_ "strconv"
)

func clear() {
	// Wipes the storage so one test does not affect another
	userlib.DatastoreClear()
	userlib.KeystoreClear()
}

func TestInit(t *testing.T) {
	clear()
	t.Log("Initialization test")

	// You can set this to false!
	userlib.SetDebugStatus(true)

	u, err := InitUser("alice", "fubar")
	if err != nil {
		// t.Error says the test fails
		t.Error("Failed to initialize user", err)
		return
	}
	// t.Log() only produces output if you run with "go test -v"
	t.Log("Got user", &u)
	// If you want to comment the line above,
	// write _ = u here to make the compiler happy
	// You probably want many more tests here.
}

func TestGetEmptyUser(t *testing.T) {
	// Test the Get Function
	clear()
	t.Log("Initialization test")

	// You can set this to false!
	userlib.SetDebugStatus(true)

	// t.Log() only produces output if you run with "go test -v"
	// If you want to comment the line above,
	// write _ = u here to make the compiler happy
	// You probably want many more tests here.
	u1, _ := GetUser("mike", "fubar")
	
	t.Log("Check GetUser", &u1)
	// t.Log("Username?", u1.Username)

}

func TestStorage(t *testing.T) {
	clear()
	u, err := InitUser("alice", "fubar")
	if err != nil {
		t.Error("Failed to initialize user", err)
		return
	}

	v := []byte("This is a test")
	u.StoreFile("file1", v)

	v2, err2 := u.LoadFile("file1")
	if err2 != nil {
		t.Error("Failed to upload and download", err2)
		return
	}
	if !reflect.DeepEqual(v, v2) {
		t.Error("Downloaded file is not the same", v, v2)
		return
	}
}

func TestInvalidFile(t *testing.T) {
	clear()
	u, err := InitUser("alice", "fubar")
	if err != nil {
		t.Error("Failed to initialize user", err)
		return
	}

	_, err2 := u.LoadFile("this file does not exist")
	if err2 == nil {
		t.Error("Downloaded a ninexistent file", err2)
		return
	}
}

func TestUserAppendMultipleFile(t *testing.T) {
	u, err := GetUser("alice", "fubar")
	v,_ := u.LoadFile("file1")
	if err != nil {
		t.Error("Failed to download the file from alice", err)
		return
	}
	err = u.AppendFile("file1", []byte("helloworld"))
	if err != nil {
		t.Error("Failed to append file", err)
		return
	}

	err = u.AppendFile("file1", []byte("helloworld"))
	if err != nil {
		t.Error("Failed to append file", err)
		return
	}

	v2,_ := u.LoadFile("file1")
	if err != nil {
		t.Error("Failed to download the file from alice", err)
		return
	}
	v3 := append(v, []byte("helloworld")...)
	if !reflect.DeepEqual(append(v3, []byte("helloworld")...), v2) {
		t.Error("Appending wrong", v3, v2)
		return
	}
}


func TestShare(t *testing.T) {
	clear()
	u, err := InitUser("alice", "fubar")
	if err != nil {
		t.Error("Failed to initialize user", err)
		return
	}
	u2, err2 := InitUser("bob", "foobar")
	if err2 != nil {
		t.Error("Failed to initialize bob", err2)
		return
	}

	v := []byte("This is a test")
	u.StoreFile("file1", v)
	
	var v2 []byte
	var magic_string string

	v, err = u.LoadFile("file1")
	if err != nil {
		t.Error("Failed to download the file from alice", err)
		return
	}

	magic_string, err = u.ShareFile("file1", "bob")
	if err != nil {
		t.Error("Failed to share the a file", err)
		return
	}
	err = u2.ReceiveFile("file2", "alice", magic_string)
	if err != nil {
		t.Error("Failed to receive the share message", err)
		return
	}

	v2, err = u2.LoadFile("file2")
	if err != nil {
		t.Error("Failed to download the file after sharing", err)
		return
	}
	if !reflect.DeepEqual(v, v2) {
		t.Error("Shared file is not the same", v, v2)
		return
	}

}

func TestShareandAppend(t *testing.T) {
	clear()
	userlib.SetDebugStatus(true)
	u, err := InitUser("alice", "fubar")
	if err != nil {
		t.Error("Failed to initialize user", err)
		return
	}
	u2, err2 := InitUser("bob", "foobar")
	if err2 != nil {
		t.Error("Failed to initialize bob", err2)
		return
	}

	v := []byte("This is a test")
	u.StoreFile("file1", v)
	
	

	v, err = u.LoadFile("file1")
	if err != nil {
		t.Error("Failed to download the file from alice", err)
		return
	}

	magic_string, err := u.ShareFile("file1", "bob")
	if err != nil {
		t.Error("Failed to share the a file", err)
		return
	}
	err = u2.ReceiveFile("file2", "alice", magic_string)
	if err != nil {
		t.Error("Failed to receive the share message", err)
		return
	}

	v2, err := u2.LoadFile("file2")
	if err != nil {
		t.Error("Before Bob append Failed to download the file after sharing", err)
		return
	}
	if !reflect.DeepEqual(v, v2) {
		t.Error("Shared file is not the same", v, v2)
		return
	}

	err = u2.AppendFile("file2", []byte("some random testing"))
	if err != nil {
		t.Error("BOB Failed to append file", err)
		return
	}
	v3 := append(v, []byte("some random testing")...)
	v4, err := u.LoadFile("file1")
	if err != nil {
		t.Error(" After Bob append, ALice Failed to download the file after sharing", err)
		return
	}
	if !reflect.DeepEqual(v3, v4) {
		t.Error("Alice's file is not updated after BOB's append", v3, v4)
		return
	}



}

