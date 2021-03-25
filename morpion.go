package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	PLAYER1 int = 1
	PLAYER2 int = 2
	xMap    int = 3
	yMap    int = 3
	WIN     int = 1
	NULL    int = 0
)

var (
	gameMap        [xMap][yMap]int
	tempGameMap    [xMap][yMap]int
	gameMapSave    [xMap * yMap][xMap][yMap]int
	userCharacter  = [2]int{11, 22}
	pos            int
	err            int
	partyEnd       bool = true
	playerTurn     int  = PLAYER1
	numberOfRounds int  = 0
)

func main() {
	gameMap = initializeMap()
	printMap(gameMap)

	for partyEnd {
		/**
		*	Here we ask & set the futur player position,
		*   We also modify the map into a temporary map using the position of the pion given
		 */
		pos = askPlayer(playerTurn)
		tempGameMap, err = modifyMap(pos, playerTurn, gameMap)

		/**
		*	The function modifyMap returns a "err" variable describing if the position asked is already used, 0 = no, 1 = yes.. if it is, we reset the loop & re ask the player
		 */
		if err == 1 {
			fmt.Println("	La case choisie est déjà utilisée, veuillez en choisir une autre.")
			continue
		}

		gameMapSave[numberOfRounds] = tempGameMap // Here we save the actual map to get a history of the game
		printMap(tempGameMap)
		gameMap = tempGameMap
		state, player := gameState(gameMap) // This function checks if someone is winning, and who
		playerTurn = changePlayerTurn(playerTurn)
		numberOfRounds++

		/**
		* If the player want, he can do a rewind of his game using the replay function
		 */
		if state == WIN {
			fmt.Println("Le joueur n °", player+1, "a gagné la partie !")
			replay(gameMapSave)
			break
		} else if numberOfRounds == 9 {
			fmt.Println("Match nul, personne n'a gagné !")
			replay(gameMapSave)
			break
		}
	}
}

/**
* This function is used to send the futur player
 */
func changePlayerTurn(actualPlayer int) int {
	if actualPlayer == PLAYER1 {
		return PLAYER2
	} else {
		return PLAYER1
	}
}

/**
* Here we check if a player won the game
 */
func gameState(gameMap [xMap][yMap]int) (int, int) {
	for x := 0; x < 2; x++ {
		//Horizontale
		for i := 0; i < 3; i++ {
			if gameMap[i][0] == userCharacter[x] && gameMap[i][1] == userCharacter[x] && gameMap[i][2] == userCharacter[x] {
				return WIN, x
			}
		}

		//Verticale
		for i := 0; i < 3; i++ {
			if gameMap[0][i] == userCharacter[x] && gameMap[1][i] == userCharacter[x] && gameMap[2][i] == userCharacter[x] {
				return WIN, x
			}
		}

		//Diagonale
		if gameMap[0][0] == userCharacter[x] && gameMap[1][1] == userCharacter[x] && gameMap[2][2] == userCharacter[x] {
			return 1, x
		} else if gameMap[0][2] == userCharacter[x] && gameMap[1][1] == userCharacter[x] && gameMap[2][0] == userCharacter[x] {
			return WIN, x
		}
	}

	// Nobody won...
	return NULL, 0
}

func modifyMap(pos int, player int, gameMap [xMap][yMap]int) ([xMap][yMap]int, int) {
	for i := 0; i < 3; i++ {
		for x := 0; x < 3; x++ {
			if gameMap[i][x] == pos {
				gameMap[i][x] = userCharacter[player-1]
				return gameMap, 0
			}
		}
	}

	return gameMap, 1
}

func printMap(gameMap [xMap][yMap]int) {
	x := 0
	for i := 0; i < 5; i++ {
		if i == 1 || i == 3 {
			fmt.Println("――┼―――┼―――")
		} else {
			fmt.Println(c(gameMap[x][0]), "|", c(gameMap[x][1]), "|", c(gameMap[x][2]))
			x++
		}
	}
}

/**
* Used to convert userCharacter into a mostly readable character
* due to int array, we cannot do this into the array... we need to do it in postprod
 */
func c(pos int) string {
	if pos == userCharacter[PLAYER1-1] {
		return "X"
	} else if pos == userCharacter[PLAYER2-1] {
		return "O"
	} else {
		return strconv.Itoa(pos)
	}
}

/**
* This function set the pattern 1 2 3 | 4 5 6 | 7 8 9 for the map
 */
func initializeMap() [3][3]int {
	var gMap [xMap][yMap]int
	x := 0
	y := 0

	for i := 1; i < 10; i++ {
		gMap[x][y] = i
		if i%3 == 0 {
			x++
			y = 0
		} else {
			y++
		}
	}

	return gMap
}

/**
* Used to see a rewind of the game
 */
func replay(gameMapSave [xMap * yMap][xMap][yMap]int) {
	for true {
		fmt.Printf("\nVoulez vous voir une rediffusion de votre match, 1/0 : ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		entry, err := strconv.Atoi(scanner.Text()) // Convert string into int

		if err != nil { // If err != nil the user is inputing a string or something that can't be converted into a int
			fmt.Println("Veuillez rentrer 1 pour oui ou 0 pour non.")
			continue
		} else if entry != 0 && entry != 1 {
			fmt.Println("Veuillez rentrer 1 pour oui ou 0 pour non.")
			continue
		}

		if entry == 0 {
			fmt.Println("A la prochaine !")
			break
		} else { // Here we print the rewind
			for i := 0; i < numberOfRounds; i++ {
				fmt.Println("")
				printMap(gameMapSave[i])
			}

			fmt.Println("\nA la prochaine !")
			break
		}
	}
}

/**
* Used for asking the futur position of the player
 */
func askPlayer(player int) int {
	var playerPos int
	for true {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("[Joueur %v] Choissisez votre case, comprise entre 1 et 9 : ", player)
		scanner.Scan()
		pos, err := strconv.Atoi(scanner.Text()) // Convert string into int

		switch {
		case err != nil: // If err != nil the user is inputing a string or something that can't be converted into a int
			fmt.Println("Veuillez entrer une position entre 1 et 9.")
			continue
		case pos < 1:
			fmt.Println("Veuillez entrer une position supérieure à 0.")
			continue
		case pos > 9:
			fmt.Println("Veuillez entrer une position inférieure ou égale à 9.")
		}

		// We check if the entry of the user is correct
		if pos >= 1 && pos <= 9 {
			playerPos = pos
			break
		}
	}

	return playerPos
}
