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

package bench

import "testing"

/*
* @author Damon
* @date   2023/7/3 9:27
 */

func BenchmarkDealW1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DealW1(200)
	}
}

func BenchmarkDealW2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DealW2(200)
	}
}
