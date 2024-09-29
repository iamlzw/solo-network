package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type evaluation struct {
	UniversityLevel   string `json:"universityLevel"`
	UniversityName    string `json:"universityName"`
	Faculty           string `json:"faculty"`
	TutorName         string `json:"tutorName"`
	Score             string `json:"score"`
	Description       string `json:"description"`
	Time              string `json:"time"`
	Counts            string `json:"counts"`
}

// Init is called when the smart contract is instantiated
func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := stub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "upload" {
		return s.upload(stub, args)
	} else if function == "query" {
		return s.query(stub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
  学校水平
  args[0] University level,for instance,985,211..
  学校名称
  args[1] University Name,for instance: Henan University
  院系
  args[2] Faculty ,for instance: Compute Science
  导师名称
  args[3] Tutor name
  评分
  args[4] score
  描述
  args[5] describe
  时间
  args[6] time
  统计
  args[7] counts
*/


func (s *SmartContract) upload(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	// Check we have a valid number of args
        fmt.Println(len(args))
	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments, expecting 8")
	}
	compositeIndexName := args[0]+"~"+args[1]+"~"+args[2]

	// Create the composite key that will allow us to query for all deltas on a particular variable
	compositeKey, compositeErr := stub.CreateCompositeKey(compositeIndexName, []string{args[3],args[4], args[5],args[6],args[7]})
	if compositeErr != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", args[3], compositeErr.Error()))
	}

	// Save the composite key index
	compositePutErr := stub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return shim.Error(fmt.Sprintf("Could not put operation for %s in the ledger: %s", args[3], compositePutErr.Error()))
	}

	return shim.Success([]byte(fmt.Sprintf("Successfully added eval")))
}

func (s *SmartContract) query(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	// Check we have a valid number of args
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments, expecting 4")
	}

	compositeIndexName := args[0]+"~"+args[1]+"~"+args[2]

	tutorName := args[3]
	// Get all deltas for the variable
	deltaResultsIterator, deltaErr := stub.GetStateByPartialCompositeKey(compositeIndexName, []string{tutorName})
	if deltaErr != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve value for %s: %s", tutorName, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return shim.Error(fmt.Sprintf("No variable by the name %s exists", tutorName))
	}

	var i int
	var evaluations []evaluation
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return shim.Error(nextErr.Error())
		}

		// Split the composite key into its component parts
		_, keyParts, splitKeyErr := stub.SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			return shim.Error(splitKeyErr.Error())
		}

		// Retrieve the delta value and operation

		score := keyParts[1]
		describe := keyParts[2]
		time := keyParts[3]
		counts := keyParts[4]
		evaluations = append(evaluations, evaluation{args[0],args[1],
			args[2],args[3],score,describe,time,counts})

	}
	evaluationsBytes,_ := json.Marshal(evaluations)

	return shim.Success(evaluationsBytes)
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

