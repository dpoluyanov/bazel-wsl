/*
   Copyright 2021 Dmitriy Poluyanov

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package bep

import (
	"bazel-wsl/bep_proto/buildeventstream"
	"bazel-wsl/utils"
	"fmt"
	"github.com/golang/protobuf/proto"
	"os"
	"strings"
)

func main() {
	var in, err = os.OpenFile("./bep.bin", os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}

	out, err := os.OpenFile("./bep.out", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	RewriteBep(in, out, out)
}

func RewriteBep(from *os.File, to *os.File, logTo *os.File) {
	data, err := os.ReadFile(from.Name())
	if err != nil {
		panic(err)
	}

	var outBuf = proto.NewBuffer(nil)

	commonHomePrefix := utils.WSLToWinPath("/home")
	for i := 0; i < len(data); {
		var size, len = proto.DecodeVarint(data[i:])

		var msgData = data[i : i+len+int(size)]

		buildEvent := &buildeventstream.BuildEvent{}
		err = proto.NewBuffer(msgData).DecodeMessage(buildEvent)
		if err != nil {
			panic(err)
		}

		var eventId = buildEvent.GetId().Id
		switch eventId.(type) {
		case *buildeventstream.BuildEventId_Workspace:
			var localExecRoot = buildEvent.GetWorkspaceInfo().GetLocalExecRoot()
			fmt.Printf("localExecRoot is %s\n", localExecRoot)
			break
		case *buildeventstream.BuildEventId_Configuration:
			var configurationId = buildEvent.GetId().GetConfiguration().GetId()
			var mnemonic = buildEvent.GetConfiguration().GetMnemonic()
			fmt.Printf("configurationId/mnemonic is %s/%s\n", configurationId, mnemonic)
			break
		case *buildeventstream.BuildEventId_NamedSet:
			//var namedSet = buildEvent.GetNamedSetOfFiles()
			// todo patch (replace) file:/// prefix (todo 0cache wsl prefix & file name once for performance)
			//fmt.Println(namedSet)
			var namedSet = buildEvent.GetNamedSetOfFiles()
			var newFiles = make([]*buildeventstream.File, 0)
			for _, file := range namedSet.Files {
				originalUri := file.GetUri()
				var path = strings.Replace(originalUri, "file:///home", commonHomePrefix, 1)
				//var newUri = utils.WSLToWinPath(path)
				path = "file://" + strings.ReplaceAll(path, "/", "\\")
				file.File = &buildeventstream.File_Uri{Uri: path}
				logTo.WriteString(fmt.Sprintf("Replaced %s to %s\n", originalUri, path))
				newFiles = append(newFiles, file)
			}
			buildEvent.GetNamedSetOfFiles().Files = newFiles
			break
		case *buildeventstream.BuildEventId_TargetCompleted:
			label := buildEvent.GetId().GetTargetCompleted().GetLabel()
			configId := buildEvent.GetId().GetTargetCompleted().GetConfiguration().GetId()
			fmt.Printf("label/configId is %s/%s\n", label, configId)
			break
		case *buildeventstream.BuildEventId_Started:
			// ignored
			break
		case *buildeventstream.BuildEventId_BuildFinished:
			// ignored
			break
		default:
			break
		}

		i += len
		i += int(size)

		// transfer part
		var err = outBuf.EncodeMessage(buildEvent)
		if err != nil {
			panic(err)
		}
	}

	to.Write(outBuf.Bytes())
}
