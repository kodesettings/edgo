/**
    Copyright (C) 2023 - 2026, edgo authors

    This program is free software; you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation; either version 2 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License along
    with this program; if not, see <https://www.gnu.org/licenses/>.
*/

#ifndef _TEST_LANGUAGES_H_
#define _TEST_LANGUAGES_H_

//------------------------------------------------------------------------

#define GO_SAMPLE "\
package main \
\
import ( \
	\"fmt\" \
	\"time\" \
) \
\
func main() { \
	start := time.Now() \
	var count = 0 \
	for i := 0; i <= 100000000; i++ { \
		count += i \
		fmt.Println(count) \
		time.Sleep(time.Millisecond * 10) \
	} \
	fmt.Println(count, \"elapsed\", time.Since(start)) \
}\
"

//------------------------------------------------------------------------

#define PYTHON_SAMPLE "\
import time \
\
print(\"starting\") \
start_time = time.time() \
\
for i in range(100000): \
    print(i) \
    time.sleep(0.01) \
\
print(\"done\") \
\
elapsed_time = time.time() - start_time \
print(\"Elapsed time:\", elapsed_time, q\"seconds\")\
"

//------------------------------------------------------------------------

#define JAVASCRIPT_SAMPLE "function hello() { console.log('hello') }"

//------------------------------------------------------------------------

#endif // _TEST_LANGUAGES_H_
