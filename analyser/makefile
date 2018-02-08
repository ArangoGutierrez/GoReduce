#
# Copyright 2017 The GoReduce Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
analyser:
	GOOS=linux go build -o bin/analyser

test:
	@GOOS=linux go build -o bin/analyser
	./bin/analyser

.PHONY: run
run:
	./bin/analyser

.PHONY: clean
clean:
		+rm -f bin/analyser
		+rm ../data/analysis/*.dat
