package checkpoint

import (
	"fmt"
	checkpoint "github.com/CheckPointSW/cp-mgmt-api-go-sdk/APIFiles"
	"math"
	"strconv"
	"time"
)

const (
	CmeApiVersion  = "v1.3.1"
	CmeApiBasePath = "cme-api"
	CmeApiPath     = CmeApiBasePath + "/" + CmeApiVersion
)

func checkIfRequestFailed(resJson map[string]interface{}) bool {

	if resJson["status-code"] != nil {
		statusCode := resJson["status-code"].(float64)
		if int(math.Round(statusCode)) != 200 {
			return true
		}
	}
	return false
}

func buildErrorMessage(resJson map[string]interface{}) string {
	errMessage := ""
	if resJson["error"] != nil {
		errorResultJson := resJson["error"].(map[string]interface{})
		if v := errorResultJson["message"]; v != nil {
			errMessage = "Message: " + v.(string)
		}
		if v := errorResultJson["details"]; v != nil {
			errMessage += ". Details: " + v.(string)
		}
		if v := errorResultJson["error-code"]; v != nil {
			errMessage += " (Error code: " + strconv.Itoa(int(math.Round(v.(float64)))) + ")"
		}
	}
	if errMessage == "" {
		errMessage = "Request failed. For more details check cme_api logger on the management server"
	}
	return errMessage
}

func cmeObjectNotFound(resJson map[string]interface{}) bool {
	NotFoundErrorCode := []int{800, 802}
	if resJson["error"] != nil {
		errorResultJson := resJson["error"].(map[string]interface{})
		if v := errorResultJson["error-code"]; v != nil {
			errorCode := int(math.Round(v.(float64)))
			for i := range NotFoundErrorCode {
				if errorCode == NotFoundErrorCode[i] {
					return true
				}
			}
		}
	}
	return false
}

func cmeWaitForReuqest(client *checkpoint.ApiClient, requestId string) error {
	url := CmeApiBasePath + "/status/" + requestId

	for res, err := client.ApiCall(url, nil, client.GetSessionID(), true, client.IsProxyUsed(), "GET"); err != nil; {
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		data := res.GetData()
		if checkIfRequestFailed(data) {
			return fmt.Errorf(buildErrorMessage(data))
		}

		requestStatus := data["result"].(map[string]interface{})["requestStatus"].(string)
		if requestStatus == "Success" {
			return nil
		} else if requestStatus != "InProgress" {
			return fmt.Errorf(err.Error())
		}

		time.Sleep(10 * time.Second)
	}

	return fmt.Errorf("Could not get request status for request ID: %s", requestId)
}
