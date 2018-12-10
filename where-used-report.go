package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var relationset = []string{}
var relationsetXML = []string{}
var relationfile *os.File
var directory = "."

// func initFile (name string ) {

// 	var err error
// 	relationfile, err = os.Create("" + name +"-relation.txt")
// 	if err != nil {
// 		log.Fatal("Cannot create file", err)
// 	}

// 	relationfile.WriteString("digraph G {" + "\n")
// }

// func finishFile() {
//     fmt.Println("FinishFile()")

// 	relationset = removeDuplicates(relationset)

// 	for v := range relationset {
// 		relationfile.WriteString(relationset[v])
// 	}
// 	// Return the new slice.

// 	relationfile.WriteString("}")
// 	relationfile.Sync()
// 	relationfile.Close()

// }

func handler(w http.ResponseWriter, r *http.Request) {

	filepath := directory + "/" + r.URL.Path[1:] // path to a single template .oet file

	if filepath == "./favicon.ico" {
		return
	}

	relationset = []string{}    // used to store the relationships between files
	relationsetXML = []string{} // used to store the relationships between files

	// to process a directory, find all the oet files within it, then iterate through the list

	// find all template files

	// for each template do

	templateID := findTemplateID(filepath)
	findParentTemplates(templateID, filepath)

	// end do

	relationset = removeDuplicates(relationset)

	// fmt.Fprintf(w, "digraph G {"+"\n")

	// for v := range relationset {
	// 	fmt.Fprintf(w, relationset[v])
	// }

	// fmt.Fprintf(w, "}")

	fmt.Fprintf(w, "<?xml-stylesheet href='GENERIC.xslt'?>")

	for v := range relationsetXML {
		fmt.Fprintf(w, relationsetXML[v])
	}

}

func main() {

	http.HandleFunc("/", handler)
	http.HandleFunc("/GENERIC.xslt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/coni/node/Dropbox/AHS/CMIO Office/Clinical Content/XSLT/GENERIC.xslt")
	})
	http.ListenAndServe(":8081", nil)

}

func findTemplateID(path string) string {
	// return the unique identifier for the template specified in the path param

	var templateID string

	log.Printf("findTemplateID - " + path)

	var _result = grepFile(path, "<id>")
	log.Printf("findTemplateID result " + _result)

	r := strings.NewReplacer("<id>", "", "</id>", "")

	templateID = r.Replace(_result)
	templateID = strings.TrimSpace(templateID)

	return templateID
}

func findParentTemplates(id string, file string) {
	// find names of templates that contain id

	if id == "" {
		log.Printf("findParentTemplates failure....no id passed in")
		return
	}

	var foundfiles = grepDir(id)

	log.Printf("findParentTemplates( " + id + ") parents = (" + foundfiles + ")")

	results := strings.Split(foundfiles, "\n")

	relationsetXML = append(relationsetXML, "<template><filename>"+file+"</filename><contained-in>")

	for i := range results {
		result := results[i]
		parts := strings.Split(result, ":")
		parent := parts[0]
		parent = strings.TrimSpace(parent)

		if parent != "" {
			fmt.Println("findParentTemplates parent - " + parent)
			storeRelationship(directory+"/"+parent, file)
			storeRelationshipXML(directory+"/"+parent, file)
			//relationsetXML= append( relationsetXML, "<filename>" + file + "</filename>"<)
			id = findTemplateID(parent)
			findParentTemplates(id, directory+"/"+parent)
		}
	}

	relationsetXML = append(relationsetXML, "</contained-in></template>")
}

func grepDir(pattern string) string {

	cmd := exec.Command("grep", "-r", "template_id=\""+pattern)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	log.Printf("grepDir finished with error: %v", err)
	stdout := outbuf.String()
	return stdout
}

func storeRelationshipXML(parent, child string) {

	if directory != "" {
		if parent == "" {
			return
		}
	}

	relationsetXML = append(relationsetXML)

}

func grepFile(file string, pattern string) string {

	cmd := exec.Command("grep", pattern, file)

	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	log.Printf("grepFile finished with error: %v", err)
	stdout := outbuf.String()
	return stdout
}

func storeRelationship(parent, child string) {

	if directory != "" {
		if parent == "" {
			return
		}
	}

	relation := "\"" + child + "\"" + " -> " + "\"" + parent + "\"" + "\n"
	relation = strings.Replace(relation, ".oet", "", -1)
	fmt.Println(relation)
	relationset = append(relationset, relation)
}

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func findAllOETfiles(path string) string {

	cmd := exec.Command("find", path)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	log.Printf("Running command and waiting for it to finish...")
	err := cmd.Run()
	log.Printf("Command finished with error: %v", err)
	stdout := outbuf.String()
	return stdout

}

func processAllOETfiles(allfileslist string) {

	results := strings.Split(allfileslist, "\n")

	for i := range results {
		aFile := results[i]
		aFile = strings.TrimSpace(aFile)
		if strings.HasSuffix(aFile, ".oet") {
			templateID := findTemplateID(aFile)
			log.Printf(templateID)
			findParentTemplates(templateID, aFile)
		}
	}

}
