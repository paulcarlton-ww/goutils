//go:generate mockgen -destination=../mocks/logr/mockLogger.go -package=mocks github.com/go-logr/logr Logger
package kubectl

// Copied from Kraan - https://github.com/fidelity/kraan
/*
	Copyright 2020 The Kraan contributors.

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

var (
	KubectlCmd        = kubectlCmd
	KustomizeCmd      = kustomizeCmd
	NewCommandFactory = newCommandFactory
	GetFactoryPath    = CommandFactory.getPath
	GetLogger         = CommandFactory.getLogger
	GetExecProvider   = CommandFactory.getExecProvider
	GetCommandPath    = Command.getPath
	GetSubCmd         = Command.getSubCmd
	GetArgs           = Command.getArgs
	IsJSONOutput      = Command.isJSONOutput
	AsCommandString   = Command.asString
)

func SetKubectlCmd(command string) {
	kubectlCmd = command
}

func SetNewExecProviderFunc(newFunc func() ExecProvider) {
	newExecProviderFunc = newFunc
}

func SetNewTempDirProviderFunc(newFunc func() (string, error)) {
	tempDirProviderFunc = newFunc
}
