package configuration

import (
	"errors"
	"fmt"
	"os"
	"os/user"

	"golang.org/x/sys/windows/registry"
)

func CurrentPath() string {
	current_exe, _ := os.Executable()
	return current_exe
}

func CurrentDir() string {
	current_dir, _ := os.Getwd()
	return current_dir
}

func DesktopPath() string {
	myself, error := user.Current()
	desktopdir := ""
	if error != nil {
		panic(error)
	} else {
		homedir := myself.HomeDir
		desktopdir = homedir+"/Desktop/" 
	}

	return desktopdir
}

func SaveReg(steam string) {
    k, errk := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Valve\Steam`, registry.WRITE)
	if errk != nil {
		panic(errk)
	}
	defer k.Close()


	if errk = k.SetStringValue(`AutoLoginUser`, steam); errk != nil {
		panic(errk)
	}

	if errk = k.SetDWordValue(`RememberPassword`, 1); errk != nil {
		panic(errk)
	}
	
	a , b , erro := nome("als",5)
	if erro != nil{
		panic( erro)
	}
	fmt.Println(a,b)	
}

func GetSteamPath() string {

	k, errk := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Valve\Steam`, registry.QUERY_VALUE)
	if errk != nil {
		panic(errk)
	}
	defer k.Close()
	
	s, _, err := k.GetStringValue("SteamExe")
	if err != nil {
		panic(err)
	}
	
	return s
}

func GetSteamDir() string {

	k, errk := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Valve\Steam`, registry.QUERY_VALUE)
	if errk != nil {
		panic(errk)
	}
	defer k.Close()
	
	s, _, err := k.GetStringValue("SteamPath")
	if err != nil {
		panic(err)
	}
	
	return s
}

func nome(a string , b int)(string , int, error){

	if b > 10{
		return "", 0, errors.New("FFFFF")	
	}
 return a,b, nil
}