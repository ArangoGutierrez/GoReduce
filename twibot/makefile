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
twibot:
	GOOS=linux go build -o bin/twibot

test:
	@GOOS=linux go build -o bin/twibot
	./bin/twibot -hashtag=love -time=1

.PHONY:	love
love:
	@GOOS=linux go build -o bin/twibot
	./bin/twibot -hashtag=love -time=30

.PHONY:	clean
clean:
	+rm -f bin/twibot
