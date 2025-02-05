package graphql

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

type argumentBuilder interface {
	build() string
}

type nearMediaArgumentBuilder struct {
	mediaName    string
	mediaField   string
	data         string
	dataReader   io.Reader
	hasCertainty bool
	certainty    float32
	hasDistance  bool
	distance     float32
}

func (b *nearMediaArgumentBuilder) withCertainty(certainty float32) *nearMediaArgumentBuilder {
	b.hasCertainty = true
	b.certainty = certainty
	return b
}

func (b *nearMediaArgumentBuilder) withDistance(distance float32) *nearMediaArgumentBuilder {
	b.hasDistance = true
	b.distance = distance
	return b
}

func (b *nearMediaArgumentBuilder) getContent() string {
	if b.dataReader != nil {
		content, err := io.ReadAll(b.dataReader)
		if err != nil {
			return err.Error()
		}
		return base64.StdEncoding.EncodeToString(content)
	}
	if strings.HasPrefix(b.data, "data:") {
		base64 := ";base64,"
		indx := strings.LastIndex(b.data, base64)
		return b.data[indx+len(base64):]
	}
	return b.data
}

func (b *nearMediaArgumentBuilder) build() string {
	clause := []string{}
	if content := b.getContent(); content != "" {
		clause = append(clause, fmt.Sprintf("%s: \"%s\"", b.mediaField, content))
	}
	if b.hasCertainty {
		clause = append(clause, fmt.Sprintf("certainty: %v", b.certainty))
	}
	if b.hasDistance {
		clause = append(clause, fmt.Sprintf("distance: %v", b.distance))
	}
	return fmt.Sprintf("%s:{%s}", b.mediaName, strings.Join(clause, " "))
}
