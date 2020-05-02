package s3

import (
	"encoding/json"
)

func Init(filePath string) (*S3Node, error) {
	s3Requests, err := Parse(filePath)
	if err != nil {
		return nil, err
	}

	rootNode := &S3Node{Name: "Root", Path: "/", Type: "Root", children: make(map[string]*S3Node)}

	for _, s3Request := range s3Requests {
		// fmt.Println(s3Resource)
		if s3Request.Method == "PUT" {
			rootNode.addNode(s3Request.actualPath, s3Request.Data)
		}
	}
	return rootNode, nil
}

func (n *S3Node) Json(resourcePath string) ([]byte, error) {

	nodes := make([]*S3Node, 0)
	node, ok := n.getNode(resourcePath)

	if ok {
		for _, childNode := range node.children {
			nodes = append(nodes, childNode)
		}
	}

	data, err := json.Marshal(struct {
		Name     string    `json:"name"`
		Path     string    `json:"path"`
		Children []*S3Node `json:"children"`
	}{
		Name:     node.Name,
		Path:     resourcePath,
		Children: nodes,
	})
	// data, err := json.Marshal(nodes)

	if err != nil {
		return nil, err
	}
	return data, nil
}