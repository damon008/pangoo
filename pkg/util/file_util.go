/*
 * Copyright 2023-present by Damon All Rights Reserved
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"bufio"
	"io"
	"os"
)

/*
* @author Damon
* @date   2023/5/16 19:31
 */


func ReadFile(path string) (io.Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	//defer file.Close()

	// 使用 bufio 包创建带缓冲区的 Reader 对象
	reader := bufio.NewReader(file)
	/*buf := make([]byte, 1024)
	n, err := reader.Read(buf)
	fmt.Printf("print: %s", buf[:n])*/


	/*
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			hlog.Error(err)
		}
		fmt.Printf("%s", buf[:n])
	}*/

	// 返回 Reader 对象
	return reader, nil
}
