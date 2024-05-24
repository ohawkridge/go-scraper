package main

// "fmt"

// const s1 = "£31,350 - £47,839MPS (Dependent on experience)"

// const s5 = "£99,999 - £110,867LS (Dependent on experience)"
// const s2 = "MPS/UPS plus fringe"
// const s3 = "To be discussed at interview"
// const s4 = "£21,933 - £47,839UNQ/MPS/UPS"
// "£" = 163

// func parseSalary(salaryString string) (int32, int32) {
// 	var digits []rune
// 	var top []rune
// 	var bot []rune
// 	botFound := false
// 	topFound := false
// 	for _, char := range salaryString {
// 		// fmt.Println(char)
// 		// find digits (0=48 .. 9=57)
// 		if char > 47 && char < 58 {
// 			digits = append(digits, char)
// 		}

// 		// assume digits after first £ is bottom-end
// 		if char == 163 {
// 			botFound = true
// 		}

// 		if char == 163 && botFound {
// 			topFound = true
// 		}

// 		if botFound && !topFound && char > 47 && char < 58 {
// 			bot = append(bot, char)
// 		}

// 		if topFound && char > 47 && char < 58 {
// 			top = append(top, char)
// 		}

// 		fmt.Println(botFound, topFound)
// 	}

// 	// <5 digits? salary probably not given
// 	if len(digits) < 6 {
// 		return 0, 0
// 	}

// 	m1 := fmt.Sprintf("\nAll digits found (%d):\n%v\n", len(digits), digits)
// 	fmt.Println(m1)

// 	m2 := fmt.Sprintf("\nBottom digits:\n%v\n", bot)
// 	fmt.Println(m2)

// 	m3 := fmt.Sprintf("\nTop digits:\n%v\n", top)
// 	fmt.Println(m3)

// 	return 59000, 110000
// }
