package primitives

import (
	"encoding/hex"
	"reflect"
	"testing"
)

// Helper function to decode a hex string to a byte slice.
func hexDecode(s string) []byte {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		panic(err) // Can use t.Fatalf in the test case instead of panic.
	}
	return bytes
}

func TestAES(t *testing.T) {
	tests := []struct {
		name               string
		plaintextHex       string
		keyHex             string
		expectedCiphertext string
		nonceHex           string
		additionalDataHex  string
		expectedTagHex     string
	}{
		{
			name:               "AES-128",
			plaintextHex:       "00112233445566778899aabbccddeeff",
			keyHex:             "000102030405060708090a0b0c0d0e0f",
			expectedCiphertext: "69c4e0d86a7b0430d8cdb78070b4c55a",
		},
		{
			name:               "AES-192",
			plaintextHex:       "00112233445566778899aabbccddeeff",
			keyHex:             "000102030405060708090a0b0c0d0e0f1011121314151617",
			expectedCiphertext: "dda97ca4864cdfe06eaf70a0ec0d7191",
		},
		{
			name:               "AES-256",
			plaintextHex:       "00112233445566778899aabbccddeeff",
			keyHex:             "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f",
			expectedCiphertext: "8ea2b7ca516745bfeafc49904b496089",
		},
		{
			name:               "AES with zero plaintext and key",
			plaintextHex:       "00000000000000000000000000000000",
			keyHex:             "00000000000000000000000000000000",
			expectedCiphertext: "66e94bd4ef8a2c3b884cfa59ca342b2e",
		},
		{
			name:               "AES with zero plaintext",
			plaintextHex:       "00000000000000000000000000000000",
			keyHex:             "000102030405060708090a0b0c0d0e0f",
			expectedCiphertext: "c6a13b37878f5b826f4f8162a1c8d879",
		},
		{
			name:               "AES with specific key",
			plaintextHex:       "00000000000000000000000000000000",
			keyHex:             "ad7a2bd03eac835a6f620fdcb506b345",
			expectedCiphertext: "73a23d80121de2d5a850253fcf43120e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plaintext := hexDecode(tt.plaintextHex)
			key := hexDecode(tt.keyHex)
			expectedCiphertext := hexDecode(tt.expectedCiphertext)
			ciphertext, err := AESEncrypt(plaintext, key)
			if err != nil {
				t.Errorf("AESEncrypt failed: %v", err)
			}
			if !reflect.DeepEqual(ciphertext, expectedCiphertext) {
				t.Errorf("Ciphertext mismatch:\n got: %x\nwant: %x", ciphertext, expectedCiphertext)
			}
		})
	}
}

func TestEncryptGCM(t *testing.T) {
	tests := []struct {
		name                  string
		plaintextHex          string
		additionalDataHex     string
		nonceHex              string
		keyHex                string
		expectedCiphertextHex string
		expectedTagHex        string
	}{
		{
			name:                  "Test Case 1",
			plaintextHex:          "",
			additionalDataHex:     "",
			nonceHex:              "000000000000000000000000",
			keyHex:                "00000000000000000000000000000000",
			expectedCiphertextHex: "",
			expectedTagHex:        "58e2fccefa7e3061367f1d57a4e7455a",
		},
		{
			name:                  "Test Case 2",
			plaintextHex:          "00000000000000000000000000000000",
			additionalDataHex:     "",
			nonceHex:              "000000000000000000000000",
			keyHex:                "00000000000000000000000000000000",
			expectedCiphertextHex: "0388dace60b6a392f328c2b971b2fe78",
			expectedTagHex:        "ab6e47d42cec13bdf53a67b21257bddf",
		},
		{
			name:                  "Test Case 3",
			plaintextHex:          "d9313225f88406e5a55909c5aff5269a86a7a9531534f7da2e4c303d8a318a721c3c0c95956809532fcf0e2449a6b525b16aedf5aa0de657ba637b391aafd255",
			additionalDataHex:     "",
			nonceHex:              "cafebabefacedbaddecaf888",
			keyHex:                "feffe9928665731c6d6a8f9467308308",
			expectedCiphertextHex: "42831ec2217774244b7221b784d0d49ce3aa212f2c02a4e035c17e2329aca12e21d514b25466931c7d8f6a5aac84aa051ba30b396a0aac973d58e091473f5985",
			expectedTagHex:        "4d5c2af327cd64a62cf35abd2ba6fab4",
		},
		{
			name:                  "Test Case 4",
			plaintextHex:          "d9313225f88406e5a55909c5aff5269a86a7a9531534f7da2e4c303d8a318a721c3c0c95956809532fcf0e2449a6b525b16aedf5aa0de657ba637b39",
			additionalDataHex:     "feedfacedeadbeeffeedfacedeadbeefabaddad2",
			nonceHex:              "cafebabefacedbaddecaf888",
			keyHex:                "feffe9928665731c6d6a8f9467308308",
			expectedCiphertextHex: "42831ec2217774244b7221b784d0d49ce3aa212f2c02a4e035c17e2329aca12e21d514b25466931c7d8f6a5aac84aa051ba30b396a0aac973d58e091",
			expectedTagHex:        "5bc94fbc3221a5db94fae95ae7121a47",
		},
		{
			name:                  "Test Case 5",
			plaintextHex:          "d9313225f88406e5a55909c5aff5269a86a7a9531534f7da2e4c303d8a318a721c3c0c95956809532fcf0e2449a6b525b16aedf5aa0de657ba637b39",
			additionalDataHex:     "feedfacedeadbeeffeedfacedeadbeefabaddad2",
			nonceHex:              "cafebabefacedbad",
			keyHex:                "feffe9928665731c6d6a8f9467308308",
			expectedCiphertextHex: "61353b4c2806934a777ff51fa22a4755699b2a714fcdc6f83766e5f97b6c742373806900e49f24b22b097544d4896b424989b5e1ebac0f07c23f4598",
			expectedTagHex:        "3612d2e79e3b0785561be14aaca2fccb",
		},
		{
			name:                  "Test Case 6",
			plaintextHex:          "d9313225f88406e5a55909c5aff5269a86a7a9531534f7da2e4c303d8a318a721c3c0c95956809532fcf0e2449a6b525b16aedf5aa0de657ba637b39",
			additionalDataHex:     "feedfacedeadbeeffeedfacedeadbeefabaddad2",
			nonceHex:              "9313225df88406e555909c5aff5269aa6a7a9538534f7da1e4c303d2a318a728c3c0c95156809539fcf0e2429a6b525416aedbf5a0de6a57a637b39b",
			keyHex:                "feffe9928665731c6d6a8f9467308308",
			expectedCiphertextHex: "8ce24998625615b603a033aca13fb894be9112a5c3a211a8ba262a3cca7e2ca701e4a9a4fba43c90ccdcb281d48c7c6fd62875d2aca417034c34aee5",
			expectedTagHex:        "619cc5aefffe0bfa462af43c1699d050",
		},
		{
			name:                  "Test Case 7",
			plaintextHex:          "",
			additionalDataHex:     "",
			nonceHex:              "000000000000000000000000",
			keyHex:                "000000000000000000000000000000000000000000000000",
			expectedCiphertextHex: "",
			expectedTagHex:        "cd33b28ac773f74ba00ed1f312572435",
		},
		{
			name:                  "Test Case 8",
			plaintextHex:          "00000000000000000000000000000000",
			additionalDataHex:     "",
			nonceHex:              "000000000000000000000000",
			keyHex:                "000000000000000000000000000000000000000000000000",
			expectedCiphertextHex: "98e7247c07f0fe411c267e4384b0f600",
			expectedTagHex:        "2ff58d80033927ab8ef4d4587514f0fb",
		},
		{
			name:                  "Test Case 9",
			plaintextHex:          "d9313225f88406e5a55909c5aff5269a86a7a9531534f7da2e4c303d8a318a721c3c0c95956809532fcf0e2449a6b525b16aedf5aa0de657ba637b391aafd255",
			additionalDataHex:     "",
			nonceHex:              "cafebabefacedbaddecaf888",
			keyHex:                "feffe9928665731c6d6a8f9467308308feffe9928665731c",
			expectedCiphertextHex: "3980ca0b3c00e841eb06fac4872a2757859e1ceaa6efd984628593b40ca1e19c7d773d00c144c525ac619d18c84a3f4718e2448b2fe324d9ccda2710acade256",
			expectedTagHex:        "9924a7c8587336bfb118024db8674a14",
		},
		{
			name:                  "Test Case 10",
			plaintextHex:          "d9313225f88406e5a55909c5aff5269a86a7a9531534f7da2e4c303d8a318a721c3c0c95956809532fcf0e2449a6b525b16aedf5aa0de657ba637b39",
			additionalDataHex:     "feedfacedeadbeeffeedfacedeadbeefabaddad2",
			nonceHex:              "cafebabefacedbaddecaf888",
			keyHex:                "feffe9928665731c6d6a8f9467308308feffe9928665731c",
			expectedCiphertextHex: "3980ca0b3c00e841eb06fac4872a2757859e1ceaa6efd984628593b40ca1e19c7d773d00c144c525ac619d18c84a3f4718e2448b2fe324d9ccda2710",
			expectedTagHex:        "2519498e80f1478f37ba55bd6d27618c",
		},
		{
			name:                  "Test Case 11",
			plaintextHex:          "d9313225f88406e5a55909c5aff5269a86a7a9531534f7da2e4c303d8a318a721c3c0c95956809532fcf0e2449a6b525b16aedf5aa0de657ba637b39",
			additionalDataHex:     "feedfacedeadbeeffeedfacedeadbeefabaddad2",
			nonceHex:              "cafebabefacedbad",
			keyHex:                "feffe9928665731c6d6a8f9467308308feffe9928665731c",
			expectedCiphertextHex: "0f10f599ae14a154ed24b36e25324db8c566632ef2bbb34f8347280fc4507057fddc29df9a471f75c66541d4d4dad1c9e93a19a58e8b473fa0f062f7",
			expectedTagHex:        "65dcc57fcf623a24094fcca40d3533f8",
		},
		{
			name:                  "Test Case 12",
			plaintextHex:          "d9313225f88406e5a55909c5aff5269a86a7a9531534f7da2e4c303d8a318a721c3c0c95956809532fcf0e2449a6b525b16aedf5aa0de657ba637b39",
			additionalDataHex:     "feedfacedeadbeeffeedfacedeadbeefabaddad2",
			nonceHex:              "9313225df88406e555909c5aff5269aa6a7a9538534f7da1e4c303d2a318a728c3c0c95156809539fcf0e2429a6b525416aedbf5a0de6a57a637b39b",
			keyHex:                "feffe9928665731c6d6a8f9467308308feffe9928665731c",
			expectedCiphertextHex: "d27e88681ce3243c4830165a8fdcf9ff1de9a1d8e6b447ef6ef7b79828666e4581e79012af34ddd9e2f037589b292db3e67c036745fa22e7e9b7373b",
			expectedTagHex:        "dcf566ff291c25bbb8568fc3d376a6d9",
		},
		{
			name:                  "Test Case 13",
			plaintextHex:          "",
			additionalDataHex:     "",
			nonceHex:              "000000000000000000000000",
			keyHex:                "0000000000000000000000000000000000000000000000000000000000000000",
			expectedCiphertextHex: "",
			expectedTagHex:        "530f8afbc74536b9a963b4f1c4cb738b",
		},
		{
			name:                  "Test Case 14",
			plaintextHex:          "00000000000000000000000000000000",
			additionalDataHex:     "",
			nonceHex:              "000000000000000000000000",
			keyHex:                "0000000000000000000000000000000000000000000000000000000000000000",
			expectedCiphertextHex: "cea7403d4d606b6e074ec5d3baf39d18",
			expectedTagHex:        "d0d1c8a799996bf0265b98b5d48ab919",
		},
		{
			name:                  "Test Case 15",
			plaintextHex:          "d9313225f88406e5a55909c5aff5269a86a7a9531534f7da2e4c303d8a318a721c3c0c95956809532fcf0e2449a6b525b16aedf5aa0de657ba637b391aafd255",
			additionalDataHex:     "",
			nonceHex:              "cafebabefacedbaddecaf888",
			keyHex:                "feffe9928665731c6d6a8f9467308308feffe9928665731c6d6a8f9467308308",
			expectedCiphertextHex: "522dc1f099567d07f47f37a32a84427d643a8cdcbfe5c0c97598a2bd2555d1aa8cb08e48590dbb3da7b08b1056828838c5f61e6393ba7a0abcc9f662898015ad",
			expectedTagHex:        "b094dac5d93471bdec1a502270e3cc6c",
		},
		{
			name:                  "Test Case 16",
			plaintextHex:          "d9313225f88406e5a55909c5aff5269a86a7a9531534f7da2e4c303d8a318a721c3c0c95956809532fcf0e2449a6b525b16aedf5aa0de657ba637b39",
			additionalDataHex:     "feedfacedeadbeeffeedfacedeadbeefabaddad2",
			nonceHex:              "cafebabefacedbaddecaf888",
			keyHex:                "feffe9928665731c6d6a8f9467308308feffe9928665731c6d6a8f9467308308",
			expectedCiphertextHex: "522dc1f099567d07f47f37a32a84427d643a8cdcbfe5c0c97598a2bd2555d1aa8cb08e48590dbb3da7b08b1056828838c5f61e6393ba7a0abcc9f662",
			expectedTagHex:        "76fc6ece0f4e1768cddf8853bb2d551b",
		},
		{
			name:                  "Test Case 17",
			plaintextHex:          "d9313225f88406e5a55909c5aff5269a86a7a9531534f7da2e4c303d8a318a721c3c0c95956809532fcf0e2449a6b525b16aedf5aa0de657ba637b39",
			additionalDataHex:     "feedfacedeadbeeffeedfacedeadbeefabaddad2",
			nonceHex:              "cafebabefacedbad",
			keyHex:                "feffe9928665731c6d6a8f9467308308feffe9928665731c6d6a8f9467308308",
			expectedCiphertextHex: "c3762df1ca787d32ae47c13bf19844cbaf1ae14d0b976afac52ff7d79bba9de0feb582d33934a4f0954cc2363bc73f7862ac430e64abe499f47c9b1f",
			expectedTagHex:        "3a337dbf46a792c45e454913fe2ea8f2",
		},
		{
			name:                  "Test Case 18",
			plaintextHex:          "d9313225f88406e5a55909c5aff5269a86a7a9531534f7da2e4c303d8a318a721c3c0c95956809532fcf0e2449a6b525b16aedf5aa0de657ba637b39",
			additionalDataHex:     "feedfacedeadbeeffeedfacedeadbeefabaddad2",
			nonceHex:              "9313225df88406e555909c5aff5269aa6a7a9538534f7da1e4c303d2a318a728c3c0c95156809539fcf0e2429a6b525416aedbf5a0de6a57a637b39b",
			keyHex:                "feffe9928665731c6d6a8f9467308308feffe9928665731c6d6a8f9467308308",
			expectedCiphertextHex: "5a8def2f0c9e53f1f75d7853659e2a20eeb2b22aafde6419a058ab4f6f746bf40fc0c3b780f244452da3ebf1c5d82cdea2418997200ef82e44ae7e3f",
			expectedTagHex:        "a44a8266ee1c8eb0c8b5d4cf5ae9f19a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plaintext := hexDecode(tt.plaintextHex)
			additionalData := hexDecode(tt.additionalDataHex)
			nonce := hexDecode(tt.nonceHex)
			key := hexDecode(tt.keyHex)
			expectedCiphertext := hexDecode(tt.expectedCiphertextHex)
			expectedTag := hexDecode(tt.expectedTagHex)

			actualCiphertext, actualTag, err := EncryptGCM(plaintext, nonce, key, additionalData)
			if err != nil {
				t.Errorf("EncryptGCM() error = %v, wantErr false", err)
				return
			}

			if !reflect.DeepEqual(actualCiphertext, expectedCiphertext) {
				t.Errorf("Ciphertext mismatch in %v:\n got: %x\nwant: %x", tt.name, actualCiphertext, expectedCiphertext)
			}

			if !reflect.DeepEqual(actualTag, expectedTag) {
				t.Errorf("Tag mismatch in %v:\n got: %x\nwant: %x", tt.name, actualTag, expectedTag)
			}
		})
	}
}

func TestGhash(t *testing.T) {
	input := hexDecode("000000000000000000000000000000000388dace60b6a392f328c2b971b2fe7800000000000000000000000000000080")
	hashSubKey := hexDecode("66e94bd4ef8a2c3b884cfa59ca342b2e")
	expected := hexDecode("f38cbb1ad69223dcc3457ae5b6b0f885")

	actual := Ghash(input, hashSubKey)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("ghash mismatch:\n got: %x\nwant: %x", actual, expected)
	}
}
