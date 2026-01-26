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
package main\n\
\n\
import (\n\
	\"fmt\"\n\
	\"time\"\n\
)\n\
\n\
func main() {\n\
	start := time.Now()\n\
	var count = 0\n\
	for i := 0; i <= 100000000; i++ {\n\
		count += i\n\
		fmt.Println(count)\n\
		time.Sleep(time.Millisecond * 10)\n\
	}\n\
	fmt.Println(count, \"elapsed\", time.Since(start))\n\
}\n\
"

//------------------------------------------------------------------------

#define PYTHON_SAMPLE "\
import time\n\
\n\
print(\"starting\")\n\
start_time = time.time()\n\
\n\
for i in range(100000):\n\
    print(i)\n\
    time.sleep(0.01)\n\
\n\
print(\"done\")\n\
\n\
elapsed_time = time.time() - start_time\n\
print(\"Elapsed time:\", elapsed_time, q\"seconds\")\n\
"

//------------------------------------------------------------------------

#define JAVASCRIPT_SAMPLE "function hello() { console.log('hello') }"

//------------------------------------------------------------------------

#define GO_TEST_SAMPLE " \
import (\n\
	\"testing\"\n\
)\n\
\n\
func simple(t *testing.T) { }\n\
\n\
func Test1(t *testing.T) {\n\
	if 1 != 2 {\n\
		t.Errorf(\"Expected\")\n\
        }\n\
}\n\
\n\
func Test2(t *testing.T) {\n\
	if 1 != 2 { t.Errorf(\"Expected\") }\n\
}\n\
"

//------------------------------------------------------------------------

#define JAVASCRIPT_TEST_SAMPLE " \
function sum(a, b) {\n\
	return a + b;\n\
}\n\
\n\
describe(\"math tests\", () => {\n\
	it(\"positive\", () => {\n\
		expect(sum(1, 1)).toBe(2);\n\
	});\n\
\n\
	it(\"negative\", () => {\n\
		expect(sum(-1, -1)).toBe(-2);\n\
	});\n\
\n\
	it(\"failed\", () => {\n\
		expect(sum(-1, -1)).toBe(2);\n\
	});\n\
});\n\
"

//------------------------------------------------------------------------

#define JAVA_TEST_SAMPLE "\
import static org.junit.jupiter.api.Assertions.assertEquals;\n\
\n\
import calc.Calculator;\n\
import org.junit.jupiter.api.Test;\n\
\n\
class MyFirstJUnitJupiterTests {\n\
	private final Calculator calculator = new Calculator();\n\
\n\
	@Test\n\
	void addition() {\n\
		System.out.println(\"test\");\n\
		assertEquals(2, calculator.add(1, 1));\n\
	}\n\
\n\
	void simpleFunction() {\n\
		System.out.println(\"hi\");\n\
	}\n\
}\n\
"

//------------------------------------------------------------------------

#define PYTHON_TEST_SAMPLE " \
import pytest\n\
from random import randint\n\
\n\
class TestYo:\n\
	def mess(self, value):\n\
		return randint(0, value)\n\
\n\
	def test_pass(self):\n\
		assert 1 == 1\n\
\n\
	def test_fail_sometimes(self):\n\
		assert 1 == self.mess(1)\n\
\n\
"

#endif // _TEST_LANGUAGES_H_
