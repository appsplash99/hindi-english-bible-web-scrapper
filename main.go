package main

////////////////////////
import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

type SupportedLanguages struct {
	English string `json:"english"`
	Hindi   string `json:"hindi"`
}

// Defining Types

type Verses struct {
	VerseCount string `json:"verse_count"`    // 1001002
	BookNum    string `json:"book_number"`    // 1
	ChapterNum string `json:"chapter_number"` // 1

	BookName SupportedLanguages `json:"book_name"` // { English: "genesis", Hindi: "xyz" }
	Verse    SupportedLanguages `json:"verse"`     // { English: "genesis", Hindi: "xyz" }

	VerseNumber int `json:"verse_number"` // 1
}

type Bible struct {
	Verses []Verses
}

type CSSSelectors struct {
	verses_selector_for_romans_chapter_11,
	verses,
	book_name,
	chapter_number string
}

type HindiEnglishBibleDetails struct {
	book_name SupportedLanguages
	allVerses SupportedLanguages
}

type BibleBookContentsData struct {
	BookName         string `json:"bookName"`
	BookNumber       string `json:"bookNumber"`
	NumberOfChapters int    `json:"numberOfChapters"`
}

// helper - function to get verses string without verse number in the starting
func getCleanedVerses(allVerses string) []string {
	// for each chapter - prepare verses array
	clean_verses := make([]string, 0)
	all_v_sentences := strings.Split(allVerses, "\n") // array of verses
	for i := 0; i < len(all_v_sentences); i++ {
		curr_verse_sentence := all_v_sentences[i]        // "02 the whole ..."
		words := strings.Split(curr_verse_sentence, " ") // ["02", "the", "whole", ....]
		// if i == 0 {
		// 	fmt.Println("words", words)
		// }
		first_word := string(words[0]) // "02"
		if _, err := strconv.ParseInt(first_word, 10, 64); err == nil {
			// removing number from the starting of a verse
			// when first_word is number -> create new verse string without first_word
			new_v := strings.Join(words[1:], " ")
			clean_verses = append(clean_verses, new_v)
		} else {
			// when first_word is not a number
			clean_verses = append(clean_verses, curr_verse_sentence)
		}
		fmt.Println("clean_verses", clean_verses)
	}
	return clean_verses
}

const JSON_FILE_DIR = "./json-data"

// helper - function to check if json file exists in a folder
func isJSONFilePresentForABookName(lower_case_book_name string) bool {
	// reading bible book contents data
	existingJSONFiles, err := os.ReadDir(JSON_FILE_DIR)
	if err != nil {
		log.Fatal("error for json-data folder", err)
	}

	// looping over existing json files
	for _, e := range existingJSONFiles {
		existingJSONfileName := strings.ToLower(e.Name()) // "book_01__genesis.json"
		if strings.HasSuffix(existingJSONfileName, ".json") {
			isJsonFileAlreadyPresent := strings.Contains(existingJSONfileName, lower_case_book_name)
			if isJsonFileAlreadyPresent {
				// when json file already present for a bookName
				return true
			}
		}
	}

	return false
}

// ////////////////////////////////////////////
// MAIN
// ////////////////////////////////////////////
func main() {

	CSS_SELECTORS := CSSSelectors{
		book_name:                             "#main > div > div > div.textHeader > h1",
		chapter_number:                        "#main > div > div > div.textHeader > p > span",
		verses:                                "#textBody > p:nth-child(3)",
		verses_selector_for_romans_chapter_11: "#textBody > p:nth-child(2)",
	}

	english_bible_base_page_link := "https://wordproject.org/bibles/kj/"
	hindi_bible_base_page_link := "https://wordproject.org/bibles/in/"

	// browser := rod.New().SlowMotion(2 * time.Second).MustConnect().NoDefaultDevice().MustIncognito()
	browser := rod.New().MustConnect().NoDefaultDevice().MustIncognito()

	// Maintaining a dummy page
	browser.MustPage("https://wikipedia.org").MustWindowNormal()

	// number_of_verses_parsed := 0

	// //////////////////////////////////////////////////
	// looping over bible contents data
	for i := 0; i < len(bibleBookContents); i++ {
		// verses from a book
		data := []Verses{}

		max_num_of_chapters := bibleBookContents[i].NumberOfChapters
		current_book_num := bibleBookContents[i].BookNumber // 01, 02, ..., 10, 11, ..., 66

		// //// TODO: Remove after testing
		// if current_book_num != "45" {
		// 	// skip all books except romans #45
		// 	fmt.Println("\nSkipping book #", current_book_num)
		// 	continue
		// }
		// fmt.Println("Starting Romans #45...")
		// fmt.Println("\n\nbook #", current_book_num, ", max ch:", max_num_of_chapters)
		// /////////

		lower_case_book_name := strings.ToLower(bibleBookContents[i].BookName)
		isFileAlreadyPresent := isJSONFilePresentForABookName(lower_case_book_name)

		if !isFileAlreadyPresent {
			// if i == 0 {}
			// looping over each chapter
			for j := 1; j <= max_num_of_chapters; j++ {

				// if j == 1 {}
				current_chapter_num := strconv.Itoa(j) // 1, 2, 3, ..., 10, 11, ..., 66
				fmt.Println("chapter #", current_chapter_num)
				// for each chapter
				url_suffix := current_book_num + "/" + current_chapter_num + ".htm"

				englishPage := browser.MustPage(english_bible_base_page_link + url_suffix).MustWindowNormal()
				hindiPage := browser.MustPage(hindi_bible_base_page_link + url_suffix).MustWindowNormal()

				// //////////////////////////////////////////////////
				// Important: This is required because of a special case for romans chapter 11
				// function to return css selector for hindi verses
				css_verses_selector_for_hindi := func() string {
					if current_book_num == "45" && j == 11 {
						// bookname romans, chapter 11
						return CSS_SELECTORS.verses_selector_for_romans_chapter_11
					}
					// default css selector for verses present as HTML element on a page
					return CSS_SELECTORS.verses
				}()
				// //////////////////////////////////////////////////

				// //////////////////////////////////////////////////
				// when book_name is Leviticus - have a proper hindi book_name
				bible_hindi_book_name := hindiPage.MustElement(CSS_SELECTORS.book_name).MustText()
				if current_book_num == "03" || lower_case_book_name == "leviticus" {
					bible_hindi_book_name = "लैव्यवस्था"
				}
				// //////////////////////////////////////////////////

				hindiEnglishBibleDetails := HindiEnglishBibleDetails{
					book_name: SupportedLanguages{
						English: englishPage.MustElement(CSS_SELECTORS.book_name).MustText(),
						Hindi:   bible_hindi_book_name,
					},
					allVerses: SupportedLanguages{
						English: englishPage.MustElement(CSS_SELECTORS.verses).MustText(),
						Hindi:   hindiPage.MustElement(css_verses_selector_for_hindi).MustText(),
					},
				}

				// ////////////////////////////////////////////////
				// for each chapter - prepare verses array
				clean_eng_verses := getCleanedVerses(hindiEnglishBibleDetails.allVerses.English)
				clean_hin_verses := getCleanedVerses(hindiEnglishBibleDetails.allVerses.Hindi)
				// ////////////////////////////////////////////////

				// ================================================== //
				// Looping over each verse to prepare verse json data //
				// ================================================== //
				for k := 0; k < len(clean_eng_verses); k++ {
					each_cleaned_english_verse := clean_eng_verses[k]
					each_cleaned_hindi_verse := clean_hin_verses[k]
					verse_num := k + 1
					// number_of_verses_parsed++
					// fmt.Println("number_of_verses_parsed", number_of_verses_parsed)
					data = append(data, Verses{
						VerseCount:  strings.Join([]string{current_book_num, "0", current_chapter_num, "0", strconv.Itoa(verse_num)}, ""),
						ChapterNum:  current_chapter_num,
						BookNum:     current_book_num,
						VerseNumber: verse_num,
						BookName: SupportedLanguages{
							English: hindiEnglishBibleDetails.book_name.English,
							Hindi:   hindiEnglishBibleDetails.book_name.Hindi,
						},
						Verse: SupportedLanguages{
							English: each_cleaned_english_verse,
							Hindi:   each_cleaned_hindi_verse,
						},
					})
				}

				//////////////////////////////////////////////////

				// fmt.Println("Sleeping for 1 seconds")
				// time.Sleep(time.Second * 1)

				fmt.Println("Closing Pages...")
				englishPage.MustClose()
				hindiPage.MustClose()
				// fmt.Println("Current Verses")
			}

			// ======================================== //
			// saving in json
			// ======================================== //
			fmt.Println("Starting to save in json file...")

			// fmt.Println("data", data)

			// ============== //
			// saving in json //
			// ============== //
			file, jsonErr := json.MarshalIndent(data, "", "  ")
			if jsonErr != nil {
				fmt.Println("Error marshaling JSON:", jsonErr)
				return
			}

			file_name_for_each_book := JSON_FILE_DIR + "/book_" + current_book_num + "__" + lower_case_book_name + ".json"
			fileErr := os.WriteFile(file_name_for_each_book, file, 0644)
			if fileErr != nil {
				fmt.Println("Error writing JSON file:", fileErr)
				return
			}
			fmt.Println("File Created Successfully! for ", lower_case_book_name)

		} else {
			fmt.Println("~~~~~~~~~~~~~~File already present for book: --> Skipping book_name", lower_case_book_name)
		}

	}

	time.Sleep(time.Hour)
}

// ======================================== //
// bible book details data
// ======================================== //
var bibleBookContents = []BibleBookContentsData{
	{
		BookName:         "Genesis",
		BookNumber:       "01",
		NumberOfChapters: 50,
	},
	{
		BookName:         "Exodus",
		BookNumber:       "02",
		NumberOfChapters: 40,
	},
	{
		BookName:         "Leviticus",
		BookNumber:       "03",
		NumberOfChapters: 27,
	},
	{
		BookName:         "Numbers",
		BookNumber:       "04",
		NumberOfChapters: 36,
	},
	{
		BookName:         "Deuteronomy",
		BookNumber:       "05",
		NumberOfChapters: 34,
	},
	{
		BookName:         "Joshua",
		BookNumber:       "06",
		NumberOfChapters: 24,
	},
	{
		BookName:         "Judges",
		BookNumber:       "07",
		NumberOfChapters: 21,
	},
	{
		BookName:         "Ruth",
		BookNumber:       "08",
		NumberOfChapters: 4,
	},
	{
		BookName:         "1 Samuel",
		BookNumber:       "09",
		NumberOfChapters: 31,
	},
	{
		BookName:         "2 Samuel",
		BookNumber:       "10",
		NumberOfChapters: 24,
	},
	{
		BookName:         "1 Kings",
		BookNumber:       "11",
		NumberOfChapters: 22,
	},
	{
		BookName:         "2 Kings",
		BookNumber:       "12",
		NumberOfChapters: 25,
	},
	{
		BookName:         "1 Chronicles",
		BookNumber:       "13",
		NumberOfChapters: 29,
	},
	{
		BookName:         "2 Chronicles",
		BookNumber:       "14",
		NumberOfChapters: 36,
	},
	{
		BookName:         "Ezra",
		BookNumber:       "15",
		NumberOfChapters: 10,
	},
	{
		BookName:         "Nehemiah",
		BookNumber:       "16",
		NumberOfChapters: 13,
	},
	{
		BookName:         "Esther",
		BookNumber:       "17",
		NumberOfChapters: 10,
	},
	{
		BookName:         "Job",
		BookNumber:       "18",
		NumberOfChapters: 42,
	},
	{
		BookName:         "Psalms",
		BookNumber:       "19",
		NumberOfChapters: 150,
	},
	{
		BookName:         "Proverbs",
		BookNumber:       "20",
		NumberOfChapters: 31,
	},
	{
		BookName:         "Ecclesiastes",
		BookNumber:       "21",
		NumberOfChapters: 12,
	},
	{
		BookName:         "Song of Solomon",
		BookNumber:       "22",
		NumberOfChapters: 8,
	},
	{
		BookName:         "Isaiah",
		BookNumber:       "23",
		NumberOfChapters: 66,
	},
	{
		BookName:         "Jeremiah",
		BookNumber:       "24",
		NumberOfChapters: 52,
	},
	{
		BookName:         "Lamentations",
		BookNumber:       "25",
		NumberOfChapters: 5,
	},
	{
		BookName:         "Ezekiel",
		BookNumber:       "26",
		NumberOfChapters: 48,
	},
	{
		BookName:         "Daniel",
		BookNumber:       "27",
		NumberOfChapters: 12,
	},
	{
		BookName:         "Hosea",
		BookNumber:       "28",
		NumberOfChapters: 14,
	},
	{
		BookName:         "Joel",
		BookNumber:       "29",
		NumberOfChapters: 3,
	},
	{
		BookName:         "Amos",
		BookNumber:       "30",
		NumberOfChapters: 9,
	},
	{
		BookName:         "Obadiah",
		BookNumber:       "31",
		NumberOfChapters: 1,
	},
	{
		BookName:         "Jonah",
		BookNumber:       "32",
		NumberOfChapters: 4,
	},
	{
		BookName:         "Micah",
		BookNumber:       "33",
		NumberOfChapters: 7,
	},
	{
		BookName:         "Nahum",
		BookNumber:       "34",
		NumberOfChapters: 3,
	},
	{
		BookName:         "Habakkuk",
		BookNumber:       "35",
		NumberOfChapters: 3,
	},
	{
		BookName:         "Zephaniah",
		BookNumber:       "36",
		NumberOfChapters: 3,
	},
	{
		BookName:         "Haggai",
		BookNumber:       "37",
		NumberOfChapters: 2,
	},
	{
		BookName:         "Zechariah",
		BookNumber:       "38",
		NumberOfChapters: 14,
	},
	{
		BookName:         "Malachi",
		BookNumber:       "39",
		NumberOfChapters: 4,
	},
	{
		BookName:         "Matthew",
		BookNumber:       "40",
		NumberOfChapters: 28,
	},
	{
		BookName:         "Mark",
		BookNumber:       "41",
		NumberOfChapters: 16,
	},
	{
		BookName:         "Luke",
		BookNumber:       "42",
		NumberOfChapters: 24,
	},
	{
		BookName:         "John",
		BookNumber:       "43",
		NumberOfChapters: 21,
	},
	{
		BookName:         "Acts",
		BookNumber:       "44",
		NumberOfChapters: 28,
	},
	{
		BookName:         "Romans",
		BookNumber:       "45",
		NumberOfChapters: 16,
	},
	{
		BookName:         "1 Corinthians",
		BookNumber:       "46",
		NumberOfChapters: 16,
	},
	{
		BookName:         "2 Corinthians",
		BookNumber:       "47",
		NumberOfChapters: 13,
	},
	{
		BookName:         "Galatians",
		BookNumber:       "48",
		NumberOfChapters: 6,
	},
	{
		BookName:         "Ephesians",
		BookNumber:       "49",
		NumberOfChapters: 6,
	},
	{
		BookName:         "Philippians",
		BookNumber:       "50",
		NumberOfChapters: 4,
	},
	{
		BookName:         "Colossians",
		BookNumber:       "51",
		NumberOfChapters: 4,
	},
	{
		BookName:         "1 Thessalonians",
		BookNumber:       "52",
		NumberOfChapters: 5,
	},
	{
		BookName:         "2 Thessalonians",
		BookNumber:       "53",
		NumberOfChapters: 3,
	},
	{
		BookName:         "1 Timothy",
		BookNumber:       "54",
		NumberOfChapters: 6,
	},
	{
		BookName:         "2 Timothy",
		BookNumber:       "55",
		NumberOfChapters: 4,
	},
	{
		BookName:         "Titus",
		BookNumber:       "56",
		NumberOfChapters: 3,
	},
	{
		BookName:         "Philemon",
		BookNumber:       "57",
		NumberOfChapters: 1,
	},
	{
		BookName:         "Hebrews",
		BookNumber:       "58",
		NumberOfChapters: 13,
	},
	{
		BookName:         "James",
		BookNumber:       "59",
		NumberOfChapters: 5,
	},
	{
		BookName:         "1 Peter",
		BookNumber:       "60",
		NumberOfChapters: 5,
	},
	{
		BookName:         "2 Peter",
		BookNumber:       "61",
		NumberOfChapters: 3,
	},
	{
		BookName:         "1 John",
		BookNumber:       "62",
		NumberOfChapters: 5,
	},
	{
		BookName:         "2 John",
		BookNumber:       "63",
		NumberOfChapters: 1,
	},
	{
		BookName:         "3 John",
		BookNumber:       "64",
		NumberOfChapters: 1,
	},
	{
		BookName:         "Jude",
		BookNumber:       "65",
		NumberOfChapters: 1,
	},
	{
		BookName:         "Revelation",
		BookNumber:       "66",
		NumberOfChapters: 22,
	},
}
