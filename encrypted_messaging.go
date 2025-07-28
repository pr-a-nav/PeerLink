package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type EncryptedMessage struct {
	SenderID    string    `json:"sender_id"`
	ReceiverID  string    `json:"receiver_id"`
	Content     string    `json:"content"`
	Timestamp   time.Time `json:"timestamp"`
	MessageType string    `json:"message_type"`
}

type UserKeyPair struct {
	UserID     string
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

type PrivateMessaging struct {
	mu         sync.Mutex
	users      map[string]*UserKeyPair
	messages   map[string][]EncryptedMessage
	connections map[string]net.Conn
}

func NewPrivateMessaging() *PrivateMessaging {
	return &PrivateMessaging{
		users:        make(map[string]*UserKeyPair),
		messages:     make(map[string][]EncryptedMessage),
		connections:  make(map[string]net.Conn),
	}
}

func (pm *PrivateMessaging) GenerateKeyPair(userID string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	
	pm.mu.Lock()
	pm.users[userID] = &UserKeyPair{
		UserID:     userID,
		PublicKey:  &privateKey.PublicKey,
		PrivateKey: privateKey,
	}
	pm.mu.Unlock()
	
	return nil
}

func (pm *PrivateMessaging) GetPublicKey(userID string) (*rsa.PublicKey, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	user, exists := pm.users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	
	return user.PublicKey, nil
}

func (pm *PrivateMessaging) EncryptMessage(message string, recipientPublicKey *rsa.PublicKey) (string, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", err
	}
	
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	
	ciphertext := make([]byte, aes.BlockSize+len(message))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(message))
	
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, recipientPublicKey, key, nil)
	if err != nil {
		return "", err
	}
	
	encryptedData := struct {
		Key  []byte `json:"key"`
		Data []byte `json:"data"`
	}{
		Key:  encryptedKey,
		Data: ciphertext,
	}
	
	jsonData, err := json.Marshal(encryptedData)
	if err != nil {
		return "", err
	}
	
	return base64.StdEncoding.EncodeToString(jsonData), nil
}

func (pm *PrivateMessaging) DecryptMessage(encryptedMessage string, userID string) (string, error) {
	pm.mu.Lock()
	user, exists := pm.users[userID]
	pm.mu.Unlock()
	
	if !exists {
		return "", fmt.Errorf("user not found")
	}
	
	jsonData, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return "", err
	}
	
	var encryptedData struct {
		Key  []byte `json:"key"`
		Data []byte `json:"data"`
	}
	
	if err := json.Unmarshal(jsonData, &encryptedData); err != nil {
		return "", err
	}
	
	key, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, user.PrivateKey, encryptedData.Key, nil)
	if err != nil {
		return "", err
	}
	
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	
	if len(encryptedData.Data) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	
	iv := encryptedData.Data[:aes.BlockSize]
	ciphertext := encryptedData.Data[aes.BlockSize:]
	
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	
	return string(ciphertext), nil
}

func (pm *PrivateMessaging) SendPrivateMessage(senderID, receiverID, content string) error {
	recipientPublicKey, err := pm.GetPublicKey(receiverID)
	if err != nil {
		return err
	}
	
	encryptedContent, err := pm.EncryptMessage(content, recipientPublicKey)
	if err != nil {
		return err
	}
	
	message := EncryptedMessage{
		SenderID:    senderID,
		ReceiverID:  receiverID,
		Content:     encryptedContent,
		Timestamp:   time.Now(),
		MessageType: "private",
	}
	
	pm.mu.Lock()
	pm.messages[receiverID] = append(pm.messages[receiverID], message)
	pm.mu.Unlock()
	
	if conn, exists := pm.connections[receiverID]; exists {
		messageJSON, _ := json.Marshal(message)
		conn.Write(append(messageJSON, '\n'))
	}
	
	return nil
}

func (pm *PrivateMessaging) GetPrivateMessages(userID string) []EncryptedMessage {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	messages := pm.messages[userID]
	pm.messages[userID] = []EncryptedMessage{}
	
	return messages
}

func (pm *PrivateMessaging) RegisterConnection(userID string, conn net.Conn) {
	pm.mu.Lock()
	pm.connections[userID] = conn
	pm.mu.Unlock()
}

func (pm *PrivateMessaging) RemoveConnection(userID string) {
	pm.mu.Lock()
	delete(pm.connections, userID)
	pm.mu.Unlock()
} 