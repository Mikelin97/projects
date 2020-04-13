package proj2

// CS 161 Project 2 Spring 2020
// You MUST NOT change what you import.  If you add ANY additional
// imports it will break the autograder. We will be very upset.

import (
	// You neet to add with
	// go get github.com/cs161-staff/userlib
	"github.com/cs161-staff/userlib"

	// Life is much easier with json:  You are
	// going to want to use this so you can easily
	// turn complex structures into strings etc...
	"encoding/json"

	// Likewise useful for debugging, etc...
	"encoding/hex"

	// UUIDs are generated right based on the cryptographic PRNG
	// so lets make life easier and use those too...
	//
	// You need to add with "go get github.com/google/uuid"
	"github.com/google/uuid"

	// Useful for debug messages, or string manipulation for datastore keys.
	"strings"

	// Want to import errors.
	"errors"

	// Optional. You can remove the "_" there, but please do not touch
	// anything else within the import bracket.
	_ "strconv"

	// if you are looking for fmt, we don't give you fmt, but you can use userlib.DebugMsg.
	// see someUsefulThings() below:
)

// This serves two purposes: 
// a) It shows you some useful primitives, and
// b) it suppresses warnings for items not being imported.
// Of course, this function can be deleted.
func someUsefulThings() {
	// Creates a random UUID
	f := uuid.New()
	userlib.DebugMsg("UUID as string:%v", f.String())

	// Example of writing over a byte of f
	f[0] = 10
	userlib.DebugMsg("UUID as string:%v", f.String())

	// takes a sequence of bytes and renders as hex
	h := hex.EncodeToString([]byte("fubar"))
	userlib.DebugMsg("The hex: %v", h)

	// Marshals data into a JSON representation
	// Will actually work with go structures as well
	d, _ := json.Marshal(f)
	userlib.DebugMsg("The json data: %v", string(d))
	var g uuid.UUID
	json.Unmarshal(d, &g)
	userlib.DebugMsg("Unmashaled data %v", g.String())

	// This creates an error type
	userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("This is an error")))

	// And a random RSA key.  In this case, ignoring the error
	// return value
	var pk userlib.PKEEncKey
        var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("Key is %v, %v", pk, sk)
}

// Helper function: Takes the first 16 bytes and
// converts it into the UUID type
func bytesToUUID(data []byte) (ret uuid.UUID) {
	for x := range ret {
		ret[x] = data[x]
	}
	return
}

// The structure definition for a user record
type User struct {
	// For debug purpose
	Username string 
	PrivateKey userlib.PKEDecKey 
	PrivateSign userlib.DSSignKey
	PasswordKey []byte 

	FileStructUUID uuid.UUID

	// You can add other fields here if you want...
	// Note for JSON to marshal/unmarshal, the fields need to
	// be public (start with a capital letter)
}

type Head struct{
	FileUUID uuid.UUID
	FileSymEncKey []byte
	FileMACKey []byte
	Root bool 

}

type File struct {
	Prev uuid.UUID 
	PrevMAC []byte
	Data []byte
}

type FileShareRecord struct {
	Sender string 
	Recipient string 
	Shared map[string]FileShareRecord

}



// This creates a user.  It will only be called once for a user
// (unless the keystore and datastore are cleared during testing purposes)

// It should store a copy of the userdata, suitably encrypted, in the
// datastore and should store the user's public key in the keystore.

// The datastore may corrupt or completely erase the stored
// information, but nobody outside should be able to get at the stored
// User data: the name used in the datastore should not be guessable
// without also knowing the password and username.

// You are not allowed to use any global storage other than the
// keystore and the datastore functions in the userlib library.

// You can assume the password has strong entropy, EXCEPT
// the attackers may possess a precomputed tables containing 
// hashes of common passwords downloaded from the internet.
func InitUser(username string, password string) (userdataptr *User, err error) {
	var userdata User
	var fileStruct map[string]Head
	var userUUID uuid.UUID
	var saltUUID uuid.UUID
	var usernameHMACKey = []byte { 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 }
	userdataptr = &userdata

	//TODO: This is a toy implementation.
	// Generate PKE and DS for the user, 
	// save public key to KeyStore 
	// save private key to User Struct  
	userdata.Username = username
	publicKey, privateKey, _ := userlib.PKEKeyGen()
	userdata.PrivateKey = privateKey
	privateSign, publicSign, _ := userlib.DSKeyGen()
	userdata.PrivateSign = privateSign
	userlib.KeystoreSet(username + "publicKey", publicKey)
	userlib.KeystoreSet(username + "publicSign", publicSign)

	pswdByte, _ := json.Marshal(password)
	salt := userlib.RandomBytes(32)
	pswdKey := userlib.Argon2Key(pswdByte, salt, 32)
	userdata.PasswordKey = pswdKey

	fileStructUUID := uuid.New()
	userdata.FileStructUUID = fileStructUUID 

	fileStruct = make(map[string]Head)
	fileStructEncKey, fileStructMACKey := userdata.FileStructKeyGen()
	userdata.StoreFileStruct(userdata.FileStructUUID, fileStruct, fileStructEncKey, fileStructMACKey)

	// fileKey, _ := userlib.HashKDF(pswdKey, []byte("File Key"))

	userStructMACKey, _ := userlib.HashKDF(pswdKey, []byte("User Struct MAC Key"))
	userStructEncKey, _ := userlib.HashKDF(pswdKey, []byte("User Struct Encryption Key"))
	marshalUser, _ := json.Marshal(userdataptr)
	iv := userlib.RandomBytes(16)
	encryptedUser := userlib.SymEnc(userStructEncKey[:32], iv, marshalUser)
	// probably needs to digitally sign the encrypted data somehow?!!?!?!
	userStructMAC, _ := userlib.HMACEval(userStructMACKey[:16], encryptedUser)
	encryptedData := append(userStructMAC[:16], encryptedUser...)

	// DEBUGGING 
	// temp1, _ := json.Marshal(encryptedData[16:])
	// temp2, _ := json.Marshal(encryptedData)

	// userlib.DebugMsg("EncryptedData before save:%v", string(temp1))
	// userlib.DebugMsg("UserStructMAC before save:%v", string(temp2))

	// hash the username to generate UUID 
	usernameByte, _ := json.Marshal(username)
	b, _ := userlib.HMACEval(usernameHMACKey, usernameByte)
	// use username hash to store user struct in datastore
	userUUID, _ = uuid.FromBytes(b[:16])
	userlib.DatastoreSet(userUUID, encryptedData)
	//use username hash to store salt in datastore
	saltUUID, _ = uuid.FromBytes(b[17:33])
	userlib.DatastoreSet(saltUUID, salt)

	return userdataptr, nil
}


// This fetches the user information from the Datastore.  It should
// fail with an error if the user/password is invalid, or if the user
// data was corrupted, or if the user can't be found.
func GetUser(username string, password string) (userdataptr *User, err error) {
	var userdata User
	var userUUID uuid.UUID
	var saltUUID uuid.UUID
	var usernameHMACKey = []byte { 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 }

	// turn username & password from string to []byte
	usernameByte, _ := json.Marshal(username)
	b, _ := userlib.HMACEval(usernameHMACKey, usernameByte)
	userUUID, _ = uuid.FromBytes(b[:16])
	saltUUID, _ = uuid.FromBytes(b[17:33])
	//decrypt to get the user struct
	salt, _ := userlib.DatastoreGet(saltUUID)
	pswdByte, _ := json.Marshal(password)
	pswdKey := userlib.Argon2Key(pswdByte, salt, 32)
	userStructMACKey, _ := userlib.HashKDF(pswdKey, []byte("User Struct MAC Key"))
	userStructEncKey, _ := userlib.HashKDF(pswdKey, []byte("User Struct Encryption Key"))

	encryptedUser, exist := userlib.DatastoreGet(userUUID)
	if !exist  {
		return nil, errors.New(strings.ToTitle("The acquired UUID does not exist."))
	}
	if len(encryptedUser) < 16 {
		return nil, errors.New(strings.ToTitle("The acquired UUID does not contain right things."))
	}
	userStructMAC, _ := userlib.HMACEval(userStructMACKey[:16], encryptedUser[16:])

	ok := userlib.HMACEqual(userStructMAC[:16], encryptedUser[:16])
	if ok {
		marshalUser := userlib.SymDec(userStructEncKey[:32], encryptedUser[16:])
		// DEBUGGING 
		// userlib.DebugMsg("User Struct Before Unmarshal:%v", string(marshalData))
		json.Unmarshal(marshalUser, &userdata)
		userdataptr = &userdata
		return userdataptr, nil
	} else {
		return nil, errors.New(strings.ToTitle("The integrity of user struct has been compromised."))
	}
	
}

// This stores a file in the datastore.
//
// The plaintext of the filename + the plaintext and length of the filename 
// should NOT be revealed to the datastore!
func (userdata *User) StoreFile(filename string, data []byte) {
	var head Head
	var file File 
	var fileUUID uuid.UUID 


	//TODO: This is a toy implementation.
	// Get File Struct from the Datastore first
	fileStructEncKey, fileStructMACKey := userdata.FileStructKeyGen()
	fileStruct, _ := userdata.ObtainFileStruct(fileStructEncKey, fileStructMACKey)

	iv := userlib.RandomBytes(16)
	//generate a []byte with filename and username(not sure correct?!?!)
	filenameByte, _ := userlib.HMACEval(userdata.PasswordKey[:16], []byte(filename))

	// DEBUGGING
	// testHeadByte, _ := json.Marshal(string(filenameByte[:16]))
	// //headByte, _ := json.Marshal(head)
	// userlib.DebugMsg("Before storage fileStructStore as string:%v", string(testHeadByte))
	

	// thinking about hashing the filenameByte before use it as a key?? 
	// cuz if file name is too long... 

	// NEED TO ADD A CASE FOR IF THE FILE PREVIOUSLY EXISTED 
	nameEncryptionKey, _ := userlib.HashKDF(filenameByte[:16], []byte("File Name Encryption"))
	fileEncryptionKey, _ := userlib.HashKDF(filenameByte[:16], []byte("File Data Encryption"))
	fileMACKey, _ := userlib.HashKDF(filenameByte[:16], []byte("File Data MAC"))
	// headMACKey, _ := userlib.HashKDF(filenameByte[:16], []byte("Head MAC"))
	// headEncryptionKey, _ := userlib.HashKDF(filenameByte[:16], []byte("Head Encryption"))

	//encrypt the filename with the encryption key (probably need to hash somehow)
	// if performance is a problem, maybe not need to encrypt??
	encryptedFilename := userlib.SymEnc(nameEncryptionKey[:32], iv, filenameByte)
	fileUUID, _ = uuid.FromBytes(encryptedFilename[:16])


	//put the data and (maybe a bunch of other things in this???!?!) 
	file.Data = data
	//maybe need to initiate other variables for File in here? 
	fileByte, _ := json.Marshal(file)
	encryptedData := userlib.SymEnc(fileEncryptionKey[:32], iv, fileByte)
	macData, _ := userlib.HMACEval(fileMACKey[:16], encryptedData)
	toStore := append(macData[:16], encryptedData...)
	// probably need to digital sign the encrypted data too 
	head.FileMACKey = fileMACKey[:16]
	head.FileSymEncKey = fileEncryptionKey
	head.FileUUID = fileUUID
	head.Root = true 
	fileStruct[string(filenameByte[:16])] = head
	// DEBUGGING

	// userlib.DebugMsg("Before storage file as string:%v", string(fileByte))
	// toStoreByte, _ := json.Marshal(toStore[16:])
	// userlib.DebugMsg("Before storage The encryptedData as string:%v", string(toStoreByte))
	// // headByte, _ := json.Marshal(head)
	// userlib.DebugMsg("Before storage head as string:%v", string(headByte))
	// after the marshal the head file,
	// we probably need to sign it or HMAC it somehow 

	// DEBUGGING 
	// userlib.DebugMsg("Before storage file Structure as string:%v", string(fileStructByte))

	userdata.StoreFileStruct(userdata.FileStructUUID, fileStruct, fileStructEncKey, fileStructMACKey)
	userlib.DatastoreSet(fileUUID, toStore)
	//End of toy implementation

	return
}

// This adds on to an existing file.
//
// Append should be efficient, you shouldn't rewrite or reencrypt the
// existing file, but only whatever additional information and
// metadata you need.
func (userdata *User) AppendFile(filename string, data []byte) (err error) {
	var file File 



	// Get File Struct from the Datastore first
	fileStructEncKey, fileStructMACKey := userdata.FileStructKeyGen()
	fileStruct, _ := userdata.ObtainFileStruct(fileStructEncKey, fileStructMACKey)
	//getting the HEAD struct from FileStruct 
	head, filenameByte1 := userdata.ObtainHeadStruct(filename, fileStruct)
	

	// DEBUGGING 
	// headByte, _ := json.Marshal(head)
	// userlib.DebugMsg("Checking whether head exists in APPEND:%v", string(headByte))
	if head.FileSymEncKey == nil {
		return errors.New(strings.ToTitle("This is an empty head file!"))
	}

	prevEncryptedData, _ := userlib.DatastoreGet(head.FileUUID) 
	newUUID := uuid.New()
	userlib.DatastoreSet(newUUID, prevEncryptedData)

	fileEncryptionKey := head.FileSymEncKey
	fileMACKey := head.FileMACKey
	//encrypt and MAC the data comes in 
	iv := userlib.RandomBytes(16)
	file.Data = data 
	file.Prev = newUUID 
	file.PrevMAC = prevEncryptedData[:16]
	fileByte, _ := json.Marshal(file)

	// DEBUGGING 
	// userlib.DebugMsg("Checking whether the new FILE is right in APPEND:%v", string(fileByte))

	encryptedData := userlib.SymEnc(fileEncryptionKey[:32], iv, fileByte)
	macData, _ := userlib.HMACEval(fileMACKey[:16], encryptedData)

	toStore := append(macData[:16], encryptedData...)
	userlib.DatastoreSet(head.FileUUID, toStore)

	fileStruct[filenameByte1] = head 
	userdata.StoreFileStruct(userdata.FileStructUUID, fileStruct, fileStructEncKey, fileStructMACKey)

	// DEBUGGING 
	// headByte, _ := json.Marshal(head)
	// userlib.DebugMsg("The head after Bob append:%v", string(headByte))


	return
}

// This loads a file from the Datastore.
//
// It should give an error if the file is corrupted in any way.
func (userdata *User) LoadFile(filename string) (data []byte, err error) {

	//TODO: This is a toy implementation.
	
	var file File 
	var dataByte []byte

	// getting the fileStruct and Head from Datastore
	fileStructEncKey, fileStructMACKey := userdata.FileStructKeyGen()
	fileStruct, _ := userdata.ObtainFileStruct(fileStructEncKey, fileStructMACKey)
	head, _ := userdata.ObtainHeadStruct(filename, fileStruct)
	
	// DEBUGGING 
	// userlib.DebugMsg("this USER as string:%v", userdata.Username)
	// if userdata.Username == "alice" {
	// 	userlib.DebugMsg("this required file as String:%v", filename)
	// 	headByte, _ := json.Marshal(head)
	// 	userlib.DebugMsg("This head Before and After Append:%v", string(headByte))
	// }
	// thisUserByte, _ := json.Marshal(userdata.Username)
	// //testFileStruct, _ := json.Marshal(fileStruct)
	// userlib.DebugMsg("this USER as string:%v", string(thisUserByte))

	// userlib.DebugMsg("After storage file Structure as string:%v", string(fileStructByte))
	// var fileStruct1 map[string]Head 
	// json.Unmarshal(testFileStruct, &fileStruct1)
	// testFileStruct1, _ := json.Marshal(fileStruct1)
	// userlib.DebugMsg("After storage NEW NEW file Structure as string:%v", string(testFileStruct1))

	// headByte, _ := json.Marshal(head)
	// userlib.DebugMsg("After storage HEAD as string:%v", string(headByte))

	// add if else statement to make sure head is found. 
	// same thing as store file function, probably needs to hash or something 
	if head.FileSymEncKey == nil {
		return nil, errors.New(strings.ToTitle("The required filename does not exist!"))
	}
	
	//get the fileUUID, fileMAC from HEAD struct, and obtain the file from datastore
	fileMACKey := head.FileMACKey
	fileEncryptionKey := head.FileSymEncKey
	fileUUID := head.FileUUID

	//getting the file from DataStore
	toReceive, exist := userlib.DatastoreGet(fileUUID)
	if !exist {
		return nil, errors.New(strings.ToTitle("The acquired UUID does not exist."))
	}
	if len(toReceive) < 16 {
		return nil, errors.New(strings.ToTitle("The acquired UUID does not contain right things."))
	}

	downloadedFileMAC, _ := userlib.HMACEval(fileMACKey[:16], toReceive[16:])
	

	// DEBUGGING 
	// toReceiveByte, _ := json.Marshal(toReceive[16:])
	// userlib.DebugMsg("After storage EncryptedData as string:%v", string(toReceiveByte))

	ok := userlib.HMACEqual(downloadedFileMAC[:16], toReceive[:16])
	if ok {
		// This part needs to be modified, because file is stored differently.
		

		fileByte := userlib.SymDec(fileEncryptionKey[:32], toReceive[16:])
		json.Unmarshal(fileByte, &file) 
		dataByte = append(dataByte, file.Data...)
		prevMAC := file.PrevMAC 
		// DEBUGGING 
		// prevMACByte, _ := json.Marshal(prevMAC)
		// userlib.DebugMsg("After storage Prev KeyMAC as string:%v", string(prevMACByte))
		// userlib.DebugMsg("After storage file as string:%v", string(fileByte))

		for len(prevMAC) > 0 {

			// DEBUGGING 
			// prevMACByte, _ := json.Marshal(prevMAC)
			// userlib.DebugMsg("After storage Prev KeyMAC as string:%v", string(prevMACByte))

			prevEncryptedData, _ := userlib.DatastoreGet(file.Prev)
			if len(prevEncryptedData) < 16 {
				return nil, errors.New(strings.ToTitle("This is not the Prev FIle we are looking for."))
			}
			fileByte = userlib.SymDec(fileEncryptionKey[:32], prevEncryptedData[16:])

			// DEBUGGING 
			// userlib.DebugMsg("The file in ITERATION as string:%v", string(fileByte))
			json.Unmarshal(fileByte, &file)
			PrevData := file.Data 
			dataByte = append(PrevData, dataByte...)
			prevMAC = file.PrevMAC

		}
		return dataByte, nil 

	} else {
		return nil, errors.New(strings.ToTitle("File's integrity is not guaranteed!"))
	}

	
}

// This creates a sharing record, which is a key pointing to something
// in the datastore to share with the recipient.

// This enables the recipient to access the encrypted file as well
// for reading/appending.

// Note that neither the recipient NOR the datastore should gain any
// information about what the sender calls the file.  Only the
// recipient can access the sharing record, and only the recipient
// should be able to know the sender.
func (userdata *User) ShareFile(filename string, recipient string) (
	magic_string string, err error) {
	var rPublicKey userlib.PKEEncKey


	fileStructEncKey, fileStructMACKey := userdata.FileStructKeyGen()
	fileStruct, _ := userdata.ObtainFileStruct(fileStructEncKey, fileStructMACKey)
	head, _ := userdata.ObtainHeadStruct(filename, fileStruct)

	headByte, _ := json.Marshal(head)
	//Store Head to Datastore and encrypt, share the sym key 
	headSymKey := userlib.RandomBytes(32)
	headSymIv := userlib.RandomBytes(16)
	headStore := userlib.SymEnc(headSymKey, headSymIv, headByte)
	headUUID := uuid.New()
	userlib.DatastoreSet(headUUID, headStore)
	headUUIDByte, _ := json.Marshal(headUUID)
	rPublicKey, _ = userlib.KeystoreGet(recipient + "publicKey")
	encryptedHeadUUID, err := userlib.PKEEnc(rPublicKey, headUUIDByte) 
	encryptedHeadSymKey, err := userlib.PKEEnc(rPublicKey, headSymKey)
	encryptedHead := append(encryptedHeadUUID, encryptedHeadSymKey...)

	// DEBUGGING
	// h1 := hex.EncodeToString(encryptedHeadUUID)
	// userlib.DebugMsg("The UUID hex: %v", h1)
	// h := hex.EncodeToString(encryptedHead)
	// userlib.DebugMsg("The hex: %v", h[:1025])
	// userlib.DebugMsg("BEFORE: EncryptedHead as string:%v", string(encryptedHead))


	// DEBUGGING 
	//encryptedHeadUUIDByte, _ := json.Marshal(encryptedHeadUUID)
	// userlib.DebugMsg("Before share Encrypted HEADUUID BYTE as string:%v", string(encryptedHeadUUID))
	// userlib.DebugMsg("Before share HEAD as string:%v", string(headByte))
	// encryptedHeadByte, _ := json.Marshal(encryptedHead)
	// userlib.DebugMsg("Before share HEAD ENCRYPTED as string:%v", string(encryptedHeadByte[:500]))


	signature, _ := userlib.DSSign(userdata.PrivateSign, encryptedHead)
	magicByte := append(signature, encryptedHead...)

	// DEBUGGING
	// h := hex.EncodeToString(magicByte)
	// userlib.DebugMsg("The hex: %v", h[:256])

	return hex.EncodeToString(magicByte), err


}




// Note recipient's filename can be different from the sender's filename.
// The recipient should not be able to discover the sender's view on
// what the filename even is!  However, the recipient must ensure that
// it is authentically from the sender.
func (userdata *User) ReceiveFile(filename string, sender string,
	magic_string string) error {
	var sPublicSign userlib.DSVerifyKey
	var head Head 
	var headUUID uuid.UUID 


	sPublicSign, _ = userlib.KeystoreGet(sender + "publicSign")
	magicByte, magicErr := hex.DecodeString(magic_string)
	if magicErr != nil {
		return magicErr
	}
	signature := magicByte[:256]
	encryptedHead := magicByte[256:]
	sigErr := userlib.DSVerify(sPublicSign, encryptedHead, signature)
	if sigErr != nil {
		return sigErr
	}
	// DEBUGGING 
	// userlib.DebugMsg("AFTER: EncryptedHead as string:%v", string(encryptedHead))
	// decrypt the HEAD with private key 
	
	headUUIDByte := encryptedHead[:256]
	headSymKey := encryptedHead[256:]
	headUUIDDec, err := userlib.PKEDec(userdata.PrivateKey, headUUIDByte)
	if headUUIDDec == nil {
		return errors.New(strings.ToTitle("What the heck is going on!"))
	}
	headSymKeyDec, err := userlib.PKEDec(userdata.PrivateKey, headSymKey)
	if headSymKeyDec == nil {
		return err
	}
	json.Unmarshal(headUUIDDec, &headUUID)
	headStore, _ := userlib.DatastoreGet(headUUID)
	headByte := userlib.SymDec(headSymKeyDec, headStore)
	json.Unmarshal(headByte, &head)
	head.Root = false 

	// // DEBUGGING 
	// headByteByte, _ := json.Marshal(head)
	// userlib.DebugMsg("After share HEADBYTEBYTE as string:%v", string(headByteByte))
	//userlib.DebugMsg("After share HEADUUID:%v", string(headSymKeyDec[:38]))
	

	// getting the fileStruct from Datastore
	fileStructEncKey, fileStructMACKey := userdata.FileStructKeyGen()
	fileStruct, _ := userdata.ObtainFileStruct(fileStructEncKey, fileStructMACKey)

	filenameByte, _ := userlib.HMACEval(userdata.PasswordKey[:16], []byte(filename))
	fileStruct[string(filenameByte[:16])] = head

	userdata.StoreFileStruct(userdata.FileStructUUID, fileStruct, fileStructEncKey, fileStructMACKey)

	// DEBUGGING 
	// userlib.DebugMsg("Before storage file Structure as string:%v", string(fileStructByte))

	return nil
}

// Removes target user's access.
func (userdata *User) RevokeFile(filename string, target_username string) (err error) {
	
	// fileStruct, fileStructUUID := userdata.ObtainFileStruct() 
	// head, filenameByte1 := userdata.ObtainHeadStruct(filename, fileStruct)

	

	return
}

// Helper Methods 
func (userdata *User) FileStructKeyGen() (fileStructEncKey []byte, fileStructMACKey []byte) {
	// getting the fileStruct from Datastore
	fileStructEncKey, _ = userlib.HashKDF(userdata.PasswordKey, []byte("File Struct Encryption"))
	fileStructMACKey, _ = userlib.HashKDF(userdata.PasswordKey, []byte("File Struct MAC"))
	return fileStructEncKey, fileStructMACKey
}

func (userdata *User) StoreFileStruct(fileStructUUID uuid.UUID, fileStruct map[string]Head, fileStructEncKey []byte, fileStructMACKey []byte) (err error) {	
	// getting the fileStruct from Datastore
	iv := userlib.RandomBytes(16)
	fileStructByte, _ := json.Marshal(fileStruct)
	encryptedFileStruct := userlib.SymEnc(fileStructEncKey[:32], iv, fileStructByte)
	fileStructMAC, _ := userlib.HMACEval(fileStructMACKey[:16], encryptedFileStruct)
	toStore := append(fileStructMAC[:16], encryptedFileStruct...)
	userlib.DatastoreSet(fileStructUUID, toStore)
	return nil 
}

func (userdata *User) ObtainFileStruct(fileStructEncKey []byte, fileStructMACKey []byte) (fileStruct map[string]Head, err error) {
	// getting the fileStruct from Datastore
	toReceive, _ := userlib.DatastoreGet(userdata.FileStructUUID)
	encryptedFileStruct := toReceive[16:]
	receivedFileStructMAC, _ := userlib.HMACEval(fileStructMACKey[:16], encryptedFileStruct)
	ok := userlib.HMACEqual(receivedFileStructMAC[:16], toReceive[:16])
	if ok {
		fileStructByte := userlib.SymDec(fileStructEncKey[:32], encryptedFileStruct)
		json.Unmarshal(fileStructByte, &fileStruct)
		return fileStruct, nil
	} else {
		return nil, errors.New(strings.ToTitle("Unable to obtain the correct File Struct!"))
	}
}

func (userdata *User) ObtainHeadStruct(filename string, fileStruct map[string]Head) (head Head, filenameByte1 string) {	
	//getting the HEAD struct from FileStruct 
	filenameByte, _ := userlib.HMACEval(userdata.PasswordKey[:16], []byte(filename))
	testFilenameByte, _ := json.Marshal(string(filenameByte[:16]))
	json.Unmarshal(testFilenameByte, &filenameByte1)
	head = fileStruct[filenameByte1]
	return head, filenameByte1
}



















