package s3

import (
	"encoding/json"
	"fmt"
	"os"
)

type S3Node struct {
	Name       string
	BucketName string
	Path       string
	Type       string
	Data       string
	children   []*S3Node
}

func (n *S3Node) Print() {
	fmt.Println(fmt.Sprintf(`
	Name: %s
	Type: %s
	Path: %s
`, n.Name, n.Type, n.Path))
	if len(n.children) > 0 {
		for _, childNode := range n.children {
			childNode.Print()
		}
	}
}

func (n *S3Node) getNode(path string) (*S3Node, bool) {
	if n.Path == path {
		return n, true
	}

	for _, childNode := range n.children {
		if v, ok := childNode.getNode(path); ok {
			return v, ok
		}
	}
	return nil, false
}

func (n *S3Node) addNode(resource *S3Resource) {
	if n.Name == "Root" && len(resource.Path) == 1 {
		// adding bucket node.
		// TODO: need to check if it exists.
		bucketNode := &S3Node{
			Name:       resource.Name,
			BucketName: resource.Name,
			Path:       resource.Path[0],
			Type:       "Bucket",
		}
		n.children = append(n.children, bucketNode)
		return
	}

	if n.Name != "Root" && len(resource.Path) == 1 {
		fileNode := &S3Node{
			BucketName: n.BucketName,
			Name:       resource.Path[0],
			Type:       "File",
			Path:       fmt.Sprintf("%s/%s", n.Path, resource.Path[0]),
		}
		n.children = append(n.children, fileNode)
		return
	}

	for i := 0; i < len(n.children); i++ {
		name := n.children[i].Name
		if name == resource.Path[0] {
			resource.Path = resource.Path[1:]
			n.children[i].addNode(resource)
			return
		}
	}

	// definitely a nested resource.
	// create file node.
	dirNode := &S3Node{
		Name: resource.Path[0],
		Path: fmt.Sprintf("%s/%s", n.Path, resource.Path[0]),
		Type: "Directory",
	}
	resource.Path = resource.Path[1:]
	dirNode.addNode(resource)
	n.children = append(n.children, dirNode)
}

func (n *S3Node) GetNodesAtPath(path string) {
	nodes := make([]*S3Node, 0)
	// get nodes at a certain path.
	if v, ok := n.getNode(path); ok {
		fmt.Println("Parent: ", v)
		for _, node := range v.children {
			nodes = append(nodes, node)
		}
	}

	data, err := json.Marshal(nodes)
	if err != nil {
		panic(err)
	}

	file, _ := os.Create("./data_2.json")

	defer file.Close()

	file.Write(data)
}
