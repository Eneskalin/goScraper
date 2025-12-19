package ui



import (
	"github.com/pterm/pterm"
)



func PrintAscii(){
logo := `
_________     ________                                        
__  ____/_______  ___/__________________ _____________________
_  / __ _  __ \____ \_  ___/_  ___/  __ \/__  __ \  _ \_  ___/
/ /_/ / / /_/ /___/ // /__ _  /   / /_/ /__  /_/ /  __/  /    
\____/  \____//____/ \___/ /_/    \__,_/ _  .___/\___//_/     
                                         /_/                       
`
	pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgCyan)).Println(logo)
    


	
}