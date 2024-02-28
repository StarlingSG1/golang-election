package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type ElectionData struct {
	nbVotesByCandidate map[string]int
	department         string
	totalVotes         int
}

func showAllResults(nbVotesByCandidate map[string]int, nbVotesByDepartment map[string]int, departmentRanking map[string]int) {
	displayNbVotesByCandidate(nbVotesByCandidate)
	displayNbVotesByDepartment(nbVotesByDepartment)
	displayDepartmentRanking(departmentRanking)
}

func displayNbVotesByCandidate(nbVotesByCandidate map[string]int) {
	for candidate, votes := range nbVotesByCandidate {
		fmt.Println("Candidate:", candidate, "Votes:", votes)
	}
}

func displayNbVotesByDepartment(nbVotesByDepartment map[string]int) {
	for key, votes := range nbVotesByDepartment {
		splitKey := strings.Split(key, "_")
		candidate := splitKey[0]
		department := splitKey[1]
		fmt.Println("Department:", department, "Candidate:", candidate, "Votes:", votes)
	}
}

func displayDepartmentRanking(departmentRanking map[string]int) {
	ranking := make([]string, 0, len(departmentRanking))
	for department := range departmentRanking {
		ranking = append(ranking, department)
	}
	sort.SliceStable(ranking, func(i, j int) bool {
		return departmentRanking[ranking[i]] > departmentRanking[ranking[j]]
	})
	for i, department := range ranking {
		fmt.Println("#", i+1, ":", department)
	}
}

func readRow(row string) ElectionData {
	splitRow := strings.Split(row, ";")
	electionData := ElectionData{}
	electionData.nbVotesByCandidate = parseCandidateVotes(splitRow)
	electionData.department = splitRow[1]
	totalVotes, err := strconv.Atoi(splitRow[10])
	if err != nil {
		fmt.Println(err.Error())
	}
	electionData.totalVotes = totalVotes
	return electionData
}

func parseCandidateVotes(splitRow []string) map[string]int {
	nbVotesByCandidate := make(map[string]int)
	CANDIDATE_COLUMN_LENGTH := 7
	INDEX_OF_CANDIDATE_COLUMN_START := 23
	for i := INDEX_OF_CANDIDATE_COLUMN_START; i < len(splitRow); i += CANDIDATE_COLUMN_LENGTH {
		candidateName := splitRow[i]
		voteCountStr := splitRow[i+2]
		voteCount, _ := strconv.Atoi(voteCountStr)
		nbVotesByCandidate[candidateName] = voteCount
	}
	return nbVotesByCandidate
}

func processVotes(scanner *bufio.Scanner) (int, map[string]int, map[string]int, map[string]int) {
	IS_FIRST_ROW := true
	totalVotes := 0
	nbVotesByCandidate := make(map[string]int)
	nbVotesByDepartment := make(map[string]int)
	departmentRanking := make(map[string]int)

	for scanner.Scan() {
		if IS_FIRST_ROW {
			IS_FIRST_ROW = false
			continue
		}
		electionData := readRow(scanner.Text())
		for candidate, votes := range electionData.nbVotesByCandidate {
			candidateKey := candidate + "_" + electionData.department
			nbVotesByDepartment[candidateKey] += votes
			nbVotesByCandidate[candidate] += votes
		}
		departmentRanking[electionData.department] += electionData.totalVotes
		totalVotes += electionData.totalVotes
	}

	return totalVotes, nbVotesByCandidate, nbVotesByDepartment, departmentRanking
}

func main() {
	// Open the data file
	dataFile, err := os.Open("data.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dataFile.Close()

	scanner := bufio.NewScanner(dataFile)
	scanner.Split(bufio.ScanLines)

	totalVotes, nbVotesByCandidate, nbVotesByDepartment, departmentRanking := processVotes(scanner)

	fmt.Println("Total votes:", totalVotes)
	showAllResults(nbVotesByCandidate, nbVotesByDepartment, departmentRanking)
}
