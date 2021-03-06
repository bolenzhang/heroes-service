package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/dddengyunjie/heroes-service/blockchain"
	"github.com/dddengyunjie/heroes-service/web/common/app"
	"github.com/dddengyunjie/heroes-service/web/controllers"
	_ "github.com/dddengyunjie/heroes-service/web/routers"
	"os"
)

func initialize() (*app.Application, error) {
	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		// Network parameters
		OrdererID: "orderer.hf.chainhero.io",

		// Channel parameters
		ChannelID:     "chainhero",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/dddengyunjie/heroes-service/fixtures/artifacts/chainhero.channel.tx",

		// Chaincode parameters
		ChainCodeID:     "heroes-service",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/dddengyunjie/heroes-service/chaincode/",
		OrgAdmin:        "Admin",
		OrgName:         "org1",
		ConfigFile:      "../config.yaml",

		// User parameters
		UserName: "User1",
	}

	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return nil, err
	}
	// Close SDK
	//defer fSetup.CloseSDK()

	// Install and instantiate the chaincode

	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return nil, err
	}

	// Query the chaincode
	response, err := fSetup.QueryHello(false)
	if err != nil {
		fmt.Printf("Unable to query hello on the chaincode: %v\n", err)
	} else {
		fmt.Printf("Response from the query hello: %s\n", response)
	}
	// Invoke the chaincode
	txId, err := fSetup.InvokeHello("chainHero")
	if err != nil {
		fmt.Printf("Unable to invoke hello on the chaincode: %v\n", err)
	} else {
		fmt.Printf("Successfully invoke hello, transaction ID: %s\n", txId)
	}

	// Query again the chaincode
	response, err = fSetup.QueryHello(false)
	if err != nil {
		fmt.Printf("Unable to query hello on the chaincode: %v\n", err)
	} else {
		fmt.Printf("Response from the query hello: %s\n", response)
	}
	// Launch the web application listening
	application := &app.Application{
		Fabric: &fSetup,
	}
	return application, nil
}
func main() {
	app, err := initialize()
	if err != nil {
		return
	}
	defer app.Fabric.CloseSDK()
	beego.Router("/", &controllers.MainController{})
	beego.Router("/upload/:value:string", &controllers.UploadController{App: app}, "*:Upload")
	beego.Router("/show", &controllers.ShowController{App: app}, "get:Show;post:Show")
	beego.Router("/showHistory", &controllers.ShowController{App: app}, "get:ShowHistory;post:ShowHistory")
	beego.Router("/initMarble/:name:string/:color:string/:size:int/:owner:string", &controllers.MarbleController{App: app}, "*:InitMarble")
	beego.Router("/deleteMarble/:name:string", &controllers.MarbleController{App: app}, "*:DeleteMarble")
	beego.Router("/readMarble/:name:string", &controllers.MarbleController{App: app}, "*:ReadMarble")
	beego.Router("/transferMarble/:name:string/:newOwner:string", &controllers.MarbleController{App: app}, "*:TransferMarble")
	beego.Router("/transferMarblesBasedOnColor/:color:string/:newOwner:string", &controllers.MarbleController{App: app}, "*:TransferMarblesBasedOnColor")
	beego.Router("/queryMarblesByOwner/:owner:string", &controllers.MarbleController{App: app}, "*:QueryMarblesByOwner")
	beego.Router("/getMarblesByRange/:startKey:string/:endKey:string", &controllers.MarbleController{App: app}, "*:GetMarblesByRange")
	beego.Router("/queryMarblesWithPagination/:pageSize:int/:bookMark:string/*", &controllers.MarbleController{App: app}, "*:QueryMarblesWithPagination") //FIXME 此处路由配置不够优雅
	beego.Router("/getMarblesByRangeWithPagination/:startKey:string/:endKey:string/:pageSize:int/?:bookMark:string", &controllers.MarbleController{App: app}, "*:GetMarblesByRangeWithPagination")
	beego.Router("/queryMarbles/*", &controllers.MarbleController{App: app}, "*:QueryMarbles") //FIXME 此处路由配置不够优雅
	beego.Router("/getHistoryForMarble/:name:string", &controllers.MarbleController{App: app}, "*:GetHistoryForMarble")
	beego.Router("/testInitMarble/:name:string/:color:string/:size:int/:owner:string/:number:int", &controllers.MarbleController{App: app}, "*:TestInitMarble")

	beego.Run()
}
