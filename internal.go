package main

import (
	"fmt"
	"log"

	"github.com/chewxy/sexp"
)

func RegisterInternal(ael *Controller) {
	//ael.AddCommand(NewInternalCommand([]string{}, []string{}, "get", Get))
	//ael.AddCommand(NewInternalCommand([]string{}, []string{"ip"}, "GetExternalIP", GetIP))
	ael.AddCommand(NewInternalCommand([]string{}, []string{}, "ping", ComPing))
	//ael.AddCommand(NewInternalCommand([]string{}, []string{}, "close", ComClose))
	//ael.AddCommand(NewInternalCommand([]string{}, []string{}, "poll", ComPoll))

}

//func Get(ael *Controller, args ...string) string {
//	if len(args) == 0 {
//		return "No arguments provided"
//	}
//	var output []string
//	for _, a := range args {
//		get_command := ael.FindOutput(a)
//		if get_command.GetName() == "nil" {
//			return "N/A"
//		}
//		output = append(output, strings.TrimSpace(get_command.Run(ael, args[1:]...)))
//	}
//	result := strings.Join(output, "; ")
//	return result
//}

//func GetIP(ael *Controller, args ...string) string {
//	resp, err := http.Get("https://ifconfig.co")
//	if err != nil {
//		log.Printf("TODO: Handle error %s", err)
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Printf("TODO: Handle body read error")
//	}
//	return string(body)
//}

//func Poll(ael *Controller, args ...string) string {
//	if len(args) == 0 {
//		return ael.GetBroadcast(0)
//	} else {
//		results := make([]string, len(args))
//		for k, i := range args {
//			bid, _ := strconv.Atoi(i)
//			results[k] = strings.TrimSpace(ael.GetBroadcast(bid))
//		}
//		if len(results) <= 1 {
//			return strings.TrimSpace(results[0])
//		} else {
//			return strings.Join(results, "\n")
//		}
//	}
//}

func ComPing(ael *Controller, args sexp.Sexp) string {
	return "(DAT \"pong\")"
}

//func ComPoll(ael *Controller, args ...string) string {
//	panic("Unimplemented")
//}

func ParseHeader(head sexp.Sexp) bool {

	header := sexp.List(head.(sexp.List))
	if len(header) <= 1 {
		log.Printf("Error 0: %v\n", len(header))
		return false
	}

	log.Printf("len: %v\n value: %v\n", len(header), header)
	if header.LeafCount() < 3 || header.Head().Head() != sexp.Symbol("NEW") {
		log.Println("Error 1: no NEW")
		return false
	}
	tail := header.Tail()
	fmt.Printf("%v \n", tail)
	if tail.LeafCount() < 2 || tail.Head().Head() != sexp.Symbol("aelita") {
		log.Printf("Error 2: no Aelita")
		return false
	}
	if tail.LeafCount() < 2 || tail.Head().Tail().Head() != sexp.Symbol(PROTOV) {
		log.Printf("Error 3: Bad PROTOV")
		return false
	}
	return true
}
