package utils

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/ffandroid/proj0425/logx"
	"github.com/ffandroid/proj0425/myvideos/db"
	"github.com/ffandroid/proj0425/myvideos/repos"
	"github.com/google/uuid"
	"github.com/mr-tron/base58"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var domains []db.Domain

func init() {
	domains = repos.FindAllDomains()
}

func GetDomains() []db.Domain {
	return domains
}

// Why I have multi mongodb documents with same domainId and videoId. Doesn't that producethe same _id
func GenDocumentid(domainId int, videoId string) string {
	data := fmt.Sprintf("%d_%s", domainId, videoId)
	namespace := uuid.MustParse("5a1d2f06-1c91-4c22-9dbf-bac3b5373ad3")
	id := uuid.NewSHA1(namespace, []byte(data))
	return id.String()
}

func CompressUUID(id string) string {
	u, err := uuid.Parse(id)
	if err != nil {
		return id
	}
	return base58.Encode(u[:]) // ~22 chars, reversible
}

func DecompressUUID(short string) string {
	b, err := base58.Decode(short)
	if err != nil {
		return short
	}
	u, err := uuid.FromBytes(b)
	if err != nil {
		return short
	}
	return u.String()
}

func ExtractVideoIdX(domainId int, vurl string) string {
	vid := ""
	if domainId == 1 || domainId == 4 || domainId == 3 || domainId == 6 || domainId == 8 || domainId == 10 || domainId == 12 {
		parts := strings.Split(vurl, "/")
		if len(parts) > 4 {
			vid = parts[4]
		}
	} else if domainId == 2 || domainId == 5 || domainId == 7 || domainId == 9 || domainId == 11 {
		parts := strings.Split(vurl, "/")
		if len(parts) > 3 {
			vid = parts[3]
		}
	} else {
		vid = vurl
	}

	return vid
}
func GetDomainByID(id int) db.Domain {
	var domain db.Domain
	for _, d := range domains {
		if d.ID == id {
			domain = d
			return domain
		}
	}
	return domain
}
func ExtractVideoIdx(vurl string) string {
	vurl = strings.TrimSuffix(vurl, "/")
	parts := strings.Split(vurl, "/")
	return parts[len(parts)-1]
}
func ExtractVideoId(vurl string) string {
	if strings.Contains(vurl, "primevid.click") {
		parts := strings.Split(vurl, "#")
		if len(parts) > 1 {
			return parts[1]
		}
	}
	u, err := url.Parse(vurl)
	if err != nil {
		return ""
	}

	path := strings.Trim(u.Path, "/")
	if path == "" {
		return ""
	}

	parts := strings.Split(path, "/")
	n := len(parts)

	// Apply your rule
	if n == 1 {
		return normalizeVideoId(parts[0]) // /id
	}
	return normalizeVideoId(parts[1]) // /x/id/...anything
}

func normalizeVideoId(vid string) string {
	vid = strings.TrimPrefix(vid, "embed-")
	vid = strings.TrimSuffix(vid, ".html")
	return vid
}

func GenerateVidUrl(domainId int, videoId string) string {
	domain := GetDomainByID(domainId)
	baseUrl := "https://" + domain.Name
	if domain.AltName != "" {
		baseUrl = "https://" + domain.AltName
	}
	vidUrl := fmt.Sprintf(baseUrl+"/"+domain.VideFormat, videoId)
	return vidUrl
}

var secretKey = []byte("my32byteverysecureencryptionkey!") // Must be exactly 32 bytes

func EncryptURL(plaintext string) (string, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptURL(encrypted string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aesGCM.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	nonce := ciphertext[:aesGCM.NonceSize()]
	encryptedData := ciphertext[aesGCM.NonceSize():]

	plaintext, err := aesGCM.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func ReplaceBaseURL(original string, newBase string) string {
	origURL, err := url.Parse(original)
	if err != nil {
		return ""
	}

	newBaseURL, err := url.Parse(newBase)
	if err != nil {
		return ""
	}

	// Replace scheme and host
	origURL.Scheme = newBaseURL.Scheme
	origURL.Host = newBaseURL.Host

	return origURL.String()
}

type DuplicateGroup struct {
	ID struct {
		DomainID int    `bson:"domainId"`
		VID      string `bson:"vid"`
	} `bson:"_id"`
	IDs   []string `bson:"ids"`
	Count int      `bson:"count"`
}

// When I run it it display deleted duplicate but they still in the db
func CleanupDuplicates() error {
	ctx := context.Background()
	col := db.MyVideo2025BD.Collection("video")
	cursor, err := col.Aggregate(ctx, mongo.Pipeline{
		{{"$group", bson.D{
			{"_id", bson.D{{"domainId", "$domainId"}, {"vid", "$vid"}}},
			{"ids", bson.D{{"$push", "$_id"}}},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$match", bson.D{{"count", bson.D{{"$gt", 1}}}}}},
	})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var group DuplicateGroup
		if err := cursor.Decode(&group); err != nil {
			return err
		}

		correctID := GenDocumentid(group.ID.DomainID, strings.TrimSpace(group.ID.VID))

		// Keep first one, remove the rest
		for i, oldID := range group.IDs {
			if i == 0 {
				if oldID != correctID {
					var video db.Video
					if err := col.FindOne(ctx, bson.M{"_id": oldID}).Decode(&video); err == nil {
						video.ID = correctID
						_, err := col.InsertOne(ctx, video)
						if err == nil {
							_, _ = col.DeleteOne(ctx, bson.M{"_id": oldID})
							fmt.Printf("Fixed _id for (%d, %s)\n", group.ID.DomainID, group.ID.VID)
						}
					}
				}
			} else {
				_, _ = col.DeleteOne(ctx, bson.M{"_id": oldID})
				fmt.Printf("Deleted duplicate: %s (%d, %s)\n", oldID, group.ID.DomainID, group.ID.VID)
			}
		}
	}

	return nil
}
func CheckDuplicate() {
	// 1. Ensure no more duplicates
	ctx := context.Background()
	col := db.MyVideo2025BD.Collection("video")
	count, _ := col.CountDocuments(ctx, bson.M{})
	fmt.Println("Total documents:", count)

	cursor, _ := col.Aggregate(ctx, mongo.Pipeline{
		{{"$group", bson.D{
			{"_id", bson.D{{"domainId", "$domainId"}, {"vid", "$vid"}}},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$match", bson.D{{"count", bson.D{{"$gt", 1}}}}}},
	})
	if cursor.Next(ctx) {
		fmt.Println("⚠️ Still has duplicates")
	} else {
		fmt.Println("✅ No duplicates remain")
	}

}
func CreateFixVideoCollection() {

	ctx := context.Background()

	src := db.MyVideo2025BD.Collection("video")
	dst := db.MyVideo2025BD.Collection("video2")

	// Drop destination if exists
	_ = dst.Drop(ctx)

	// ⚙️ Step 1: Use aggregation to deduplicate by (domainId, vid)
	pipeline := mongo.Pipeline{
		{{"$sort", bson.D{{"lastUpdated", -1}}}}, // keep latest
		{{"$group", bson.D{
			{"_id", bson.D{
				{"domainId", "$domainId"},
				{"vid", "$vid"},
			}},
			{"doc", bson.D{{"$first", "$$ROOT"}}},
		}}},
	}

	cursor, err := src.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal("aggregate:", err)
	}
	defer cursor.Close(ctx)

	// ⚡ Step 2: Batch insert
	const batchSize = 1000
	var batch []interface{}
	total := 0

	for cursor.Next(ctx) {
		var result struct {
			Doc db.Video `bson:"doc"`
		}
		if err := cursor.Decode(&result); err != nil {
			logx.Println("decode:", err)
			continue
		}

		video := result.Doc
		video.ID = GenDocumentid(video.DomainID, video.VID)
		batch = append(batch, video)

		if len(batch) >= batchSize {
			if _, err := dst.InsertMany(ctx, batch); err != nil {
				logx.Println("insert batch:", err)
			}
			total += len(batch)
			fmt.Printf("Inserted %d...\n", total)
			batch = batch[:0]
		}
	}

	// Insert remaining
	if len(batch) > 0 {
		if _, err := dst.InsertMany(ctx, batch); err != nil {
			logx.Println("insert batch:", err)
		}
		total += len(batch)
	}

	fmt.Println("✅ Done. Inserted:", total)

}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
