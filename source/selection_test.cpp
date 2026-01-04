#include "selection.h"
#include <gtest/gtest.h>

TEST(SelectionTests, TestNoSelection) {
	auto text = "Hello, world!\nHow are you doing today?\nI hope you're doing well.";

	selection selection;
	selection.ssx = 0;
	selection.ssy = 0;
	selection.sex = 0;
	selection.sey = 0;

	auto got = selection.GetSelectionString(text);
	auto expected = "";
	
	EXPECT_EQ(got, expected);
}

TEST(SelectionTests, TestSingleCharacterSelection) {
	auto text = "Hello, world!\nHow are you doing today?\nI hope you're doing well.";

	selection selection;
	selection.ssx = 0;
	selection.ssy = 0;
	selection.sex = 1;
	selection.sey = 0;

	auto got = selection.GetSelectionString(text);
	auto expected = "H";
	
	EXPECT_EQ(got, expected);
}

TEST(SelectionTests, TestMultipleCharacterSelection) {
	auto text = "Hello, world!\nHow are you doing today?\nI hope you're doing well.";

	selection selection;
	selection.ssx = 0;
	selection.ssy = 0;
	selection.sex = 5;
	selection.sey = 0;

	auto got = selection.GetSelectionString(text);
	auto expected = "Hello";
	
	EXPECT_EQ(got, expected);
}

TEST(SelectionTests, TestMultipleLineSelection1) {
	auto text = "Hello, world!\nHow are you doing today?\nI hope you're doing well.";

	selection selection;
	selection.ssx = 0;
	selection.ssy = 0;
	selection.sex = 5;
	selection.sey = 1;

	auto got = selection.GetSelectionString(text);
	auto expected = "Hello, world!\nHow a";
	
	EXPECT_EQ(got, expected);
}

TEST(SelectionTests, TestMultipleLineSelection2) {
	auto text = "Hello, world!\nHow are you doing today?\nI hope you're doing well.";

	selection selection;
	selection.ssx = 0;
	selection.ssy = 0;
	selection.sex = 11;
	selection.sey = 1;

	auto got = selection.GetSelectionString(text);
	auto expected = "Hello, world!\nHow are you";
	
	EXPECT_EQ(got, expected);
}

TEST(SelectionTests, TestMultipleLineSelection3) {
	auto text = "Hello, world!\nHow are you doing today?\nI hope you're doing well.";

	selection selection;
	selection.ssx = 6;
	selection.ssy = 0;
	selection.sex = 23;
	selection.sey = 1;

	auto got = selection.GetSelectionString(text);
	auto expected = " world!\nHow are you doing today";
	
	EXPECT_EQ(got, expected);
}
