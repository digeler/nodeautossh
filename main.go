package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Extension struct {
	Properties struct {
		VMID            string `json:"vmId"`
		AvailabilitySet struct {
			ID string `json:"id"`
		} `json:"availabilitySet"`
		HardwareProfile struct {
			VMSize string `json:"vmSize"`
		} `json:"hardwareProfile"`
		StorageProfile struct {
			ImageReference struct {
				Publisher string `json:"publisher"`
				Offer     string `json:"offer"`
				Sku       string `json:"sku"`
				Version   string `json:"version"`
			} `json:"imageReference"`
			OsDisk struct {
				OsType       string `json:"osType"`
				Name         string `json:"name"`
				CreateOption string `json:"createOption"`
				Caching      string `json:"caching"`
				ManagedDisk  struct {
					StorageAccountType string `json:"storageAccountType"`
					ID                 string `json:"id"`
				} `json:"managedDisk"`
				DiskSizeGB int `json:"diskSizeGB"`
			} `json:"osDisk"`
			DataDisks []interface{} `json:"dataDisks"`
		} `json:"storageProfile"`
		OsProfile struct {
			ComputerName       string `json:"computerName"`
			AdminUsername      string `json:"adminUsername"`
			LinuxConfiguration struct {
				DisablePasswordAuthentication bool `json:"disablePasswordAuthentication"`
				SSH                           struct {
					PublicKeys []struct {
						Path    string `json:"path"`
						KeyData string `json:"keyData"`
					} `json:"publicKeys"`
				} `json:"ssh"`
				ProvisionVMAgent bool `json:"provisionVMAgent"`
			} `json:"linuxConfiguration"`
			Secrets                  []interface{} `json:"secrets"`
			AllowExtensionOperations bool          `json:"allowExtensionOperations"`
		} `json:"osProfile"`
		NetworkProfile struct {
			NetworkInterfaces []struct {
				ID string `json:"id"`
			} `json:"networkInterfaces"`
		} `json:"networkProfile"`
		ProvisioningState string `json:"provisioningState"`
	} `json:"properties"`
	Resources []struct {
		Properties struct {
			AutoUpgradeMinorVersion bool `json:"autoUpgradeMinorVersion"`
			Settings                struct {
			} `json:"settings"`
			ProvisioningState  string `json:"provisioningState"`
			Publisher          string `json:"publisher"`
			Type               string `json:"type"`
			TypeHandlerVersion string `json:"typeHandlerVersion"`
		} `json:"properties"`
		Type     string `json:"type"`
		Location string `json:"location"`
		ID       string `json:"id"`
		Name     string `json:"name"`
	} `json:"resources"`
	Type     string `json:"type"`
	Location string `json:"location"`
	Tags     struct {
		AcsengineVersion   string `json:"acsengineVersion"`
		CreationSource     string `json:"creationSource"`
		Orchestrator       string `json:"orchestrator"`
		PoolName           string `json:"poolName"`
		ResourceNameSuffix string `json:"resourceNameSuffix"`
	} `json:"tags"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Request struct {
	Location   string `json:"location"`
	Properties struct {
		Publisher          string `json:"publisher"`
		Type               string `json:"type"`
		TypeHandlerVersion string `json:"typeHandlerVersion"`
		Settings           struct {
		} `json:"settings"`
		ProtectedSettings struct {
			Username string `json:"username"`
			SSHKey   string `json:"ssh_key"`
		} `json:"protectedSettings"`
	} `json:"properties"`
}

type Vm struct {
	Value []struct {
		Properties struct {
			VMID            string `json:"vmId"`
			AvailabilitySet struct {
				ID string `json:"id"`
			} `json:"availabilitySet"`
			HardwareProfile struct {
				VMSize string `json:"vmSize"`
			} `json:"hardwareProfile"`
			StorageProfile struct {
				ImageReference struct {
					Publisher string `json:"publisher"`
					Offer     string `json:"offer"`
					Sku       string `json:"sku"`
					Version   string `json:"version"`
				} `json:"imageReference"`
				OsDisk struct {
					OsType       string `json:"osType"`
					Name         string `json:"name"`
					CreateOption string `json:"createOption"`
					Caching      string `json:"caching"`
					ManagedDisk  struct {
						StorageAccountType string `json:"storageAccountType"`
						ID                 string `json:"id"`
					} `json:"managedDisk"`
					DiskSizeGB int `json:"diskSizeGB"`
				} `json:"osDisk"`
				DataDisks []interface{} `json:"dataDisks"`
			} `json:"storageProfile"`
			OsProfile struct {
				ComputerName       string `json:"computerName"`
				AdminUsername      string `json:"adminUsername"`
				LinuxConfiguration struct {
					DisablePasswordAuthentication bool `json:"disablePasswordAuthentication"`
					SSH                           struct {
						PublicKeys []struct {
							Path    string `json:"path"`
							KeyData string `json:"keyData"`
						} `json:"publicKeys"`
					} `json:"ssh"`
					ProvisionVMAgent bool `json:"provisionVMAgent"`
				} `json:"linuxConfiguration"`
				Secrets                  []interface{} `json:"secrets"`
				AllowExtensionOperations bool          `json:"allowExtensionOperations"`
			} `json:"osProfile"`
			NetworkProfile struct {
				NetworkInterfaces []struct {
					ID string `json:"id"`
				} `json:"networkInterfaces"`
			} `json:"networkProfile"`
			ProvisioningState string `json:"provisioningState"`
		} `json:"properties"`
		Resources []struct {
			ID string `json:"id"`
		} `json:"resources"`
		Type     string `json:"type"`
		Location string `json:"location"`
		Tags     struct {
			AksEngineVersion   string `json:"aksEngineVersion"`
			CreationSource     string `json:"creationSource"`
			Orchestrator       string `json:"orchestrator"`
			PoolName           string `json:"poolName"`
			ResourceNameSuffix string `json:"resourceNameSuffix"`
		} `json:"tags"`
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"value"`
}

type Sshprop struct {
	Location   string `json:"location"`
	Properties struct {
		Publisher          string `json:"publisher"`
		Type               string `json:"type"`
		TypeHandlerVersion string `json:"typeHandlerVersion"`
		Settings           struct {
		} `json:"settings"`
		ProtectedSettings struct {
			Username string `json:"username"`
			SSHKey   string `json:"ssh_key"`
		} `json:"protectedSettings"`
	} `json:"properties"`
}

func main() {

	//get the list of vms in the rg

	args := os.Args
	if len(args) < 6 {
		fmt.Println("usage not correct \n")
		fmt.Print("usage: sub rg token ssh location")
		os.Exit(2)
	}

	subid := os.Args[1]
	rgname := os.Args[2]
	token1 := os.Args[3]
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	url := "https://management.azure.com/subscriptions/" + "/" + subid + "/resourceGroups/" + rgname + "/providers/Microsoft.Compute/virtualMachines?api-version=2018-06-01"
	//fmt.Println(url)
	//token1 := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ii1zeE1KTUxDSURXTVRQdlp5SjZ0eC1DRHh3MCIsImtpZCI6Ii1zeE1KTUxDSURXTVRQdlp5SjZ0eC1DRHh3MCJ9.eyJhdWQiOiJodHRwczovL21hbmFnZW1lbnQuY29yZS53aW5kb3dzLm5ldC8iLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC83MmY5ODhiZi04NmYxLTQxYWYtOTFhYi0yZDdjZDAxMWRiNDcvIiwiaWF0IjoxNTUwOTQxNjg0LCJuYmYiOjE1NTA5NDE2ODQsImV4cCI6MTU1MDk0NTU4NCwiX2NsYWltX25hbWVzIjp7Imdyb3VwcyI6InNyYzEifSwiX2NsYWltX3NvdXJjZXMiOnsic3JjMSI6eyJlbmRwb2ludCI6Imh0dHBzOi8vZ3JhcGgud2luZG93cy5uZXQvNzJmOTg4YmYtODZmMS00MWFmLTkxYWItMmQ3Y2QwMTFkYjQ3L3VzZXJzLzRkY2U1MGVlLWMyZGMtNGQwMy1hYjJiLTllMTVhMGMxZDdlMi9nZXRNZW1iZXJPYmplY3RzIn19LCJhY3IiOiIxIiwiYWlvIjoiQVZRQXEvOEtBQUFBS041cTJoMWhSN096RHVTZXVLcXRSSlNBTitCWVV6bW4xNEQ3RkovOXQ2MDFJVGg1UlhOL1lFbzhMYXoycHJPWERpVGJoTE00R3lGajgvYnpCclJmMllVYkFKbmhhbEJDbW02LzZGdnNvOHM9IiwiYW1yIjpbInB3ZCIsIm1mYSJdLCJhcHBpZCI6IjA0YjA3Nzk1LThkZGItNDYxYS1iYmVlLTAyZjllMWJmN2I0NiIsImFwcGlkYWNyIjoiMCIsImZhbWlseV9uYW1lIjoiR2VsZXIiLCJnaXZlbl9uYW1lIjoiRGlub3IiLCJpbl9jb3JwIjoidHJ1ZSIsImlwYWRkciI6Ijc5LjE4My4zOC4yMzYiLCJuYW1lIjoiRGlub3IgR2VsZXIiLCJvaWQiOiI0ZGNlNTBlZS1jMmRjLTRkMDMtYWIyYi05ZTE1YTBjMWQ3ZTIiLCJvbnByZW1fc2lkIjoiUy0xLTUtMjEtNzIwNTE2MDctMTc0NTc2MDAzNi0xMDkxODc5NTYtMTMyNTUwIiwicHVpZCI6IjEwMDMzRkZGODAxQzI2RTgiLCJzY3AiOiJ1c2VyX2ltcGVyc29uYXRpb24iLCJzdWIiOiJlVmQyYjlKRVlYaXFOUW5xQnltQkJITjhaMGxFdEhQWFQ4cmdBTFc5UGhvIiwidGlkIjoiNzJmOTg4YmYtODZmMS00MWFmLTkxYWItMmQ3Y2QwMTFkYjQ3IiwidW5pcXVlX25hbWUiOiJkaWdlbGVyQG1pY3Jvc29mdC5jb20iLCJ1cG4iOiJkaWdlbGVyQG1pY3Jvc29mdC5jb20iLCJ1dGkiOiJmdWtOM3ZUdEdFaVJlaEVuTTZzeEFBIiwidmVyIjoiMS4wIn0.fl6uepRWlXBPCnuvvj7Chp2fW99-1Ekca_ig8nbZVJf_HxEzTN_Pok3RYPU1EEpFCUnlUuSgMM-uy7DG1PGBRoBV6kUZpX4uZ_O2vMY1c8buKnmJtykHfK2OfKNaMvL5JTA_a6Wqw1DeQRlEbZ5ARd9p6K85d29kNGIiPE8-sQu2hcLAq9AaM61JHDJkYk7u2SS8BVzifEiYCgC8smXEXC7JKRF7DEHg6o8igX4NlW9m9edAdTakgdiYhIADlHSrU656DK_7OHKk-20r7YR_HqHjC0NooVREp-xkivZJTblsSpCrcdWpZQ_Lew2VztAet4c0GFi_kq4z6_zc2Sl4Ug"

	var bearer = "Bearer " + token1
	//fmt.Println(bearer)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	if resp.StatusCode != 200 {
		fmt.Printf("issue with token or permissions result got  " + resp.Status + " check if token didnt expired")
	}

	body, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string([]byte(body)))
	//var data = new(Vm)
	var data Vm
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}

	listvms := getvms(data)
	x := len(listvms)
	i := 0
	x = x - 1

	for i < x {

		for _, v := range listvms {

			dr, stat := Delrequest(token1, subid, rgname, v)
			if dr != nil {
				fmt.Println(dr)
			}
			fmt.Println("Delete is ok " + stat)
			pr, st, rl := Putrequest(token1, subid, rgname, v)
			if pr != nil {
				fmt.Println("error")
			}
			fmt.Println(st)
			fmt.Println(string(rl))

			err, st, lo := Checkextstat(token1, subid, rgname, v)
			if err != nil {
				fmt.Println(err)
			}
			i++
			fmt.Println(st)
			fmt.Println("itreation ", i)
			for _, vt := range lo {
				if vt != "Succeeded" {
					fmt.Println("not ready")
					fmt.Println(lo)
					//time.Sleep(5 * time.Second)
				} else {
					fmt.Println("ready to go")

					break

				}

			}

		}
	}
	for i > x {
		fmt.Println("Ready to connect sleeping....")
		time.Sleep(10 * time.Second)
	}
}

//fmt.Println(listvms)
//t := createbody()
//fmt.Println(string(t))

func getvms(s Vm) (res []string) {

	re := make([]string, 0)

	for _, nodes := range s.Value {

		re = append(re, nodes.Name)

	}

	return re
}

func createbody() (result []byte) {

	rb := new(Request)
	rb.Location = os.Args[5]
	rb.Properties.ProtectedSettings.SSHKey = os.Args[4]
	rb.Properties.ProtectedSettings.Username = "azureuser"
	rb.Properties.Publisher = "Microsoft.OSTCExtensions"
	rb.Properties.Type = "VMAccessForLinux"
	rb.Properties.TypeHandlerVersion = "1.4"
	r, _ := json.Marshal(rb)

	return r

}

func Putrequest(token string, subid string, rgname string, vmname string) (err error, st string, rl []byte) {

	url := "https://management.azure.com/subscriptions/" + subid + "/resourceGroups/" + rgname + "/providers/Microsoft.Compute/virtualMachines/" + vmname + "/extensions/enablevmaccess?api-version=2018-10-01"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var bearer = "Bearer " + token
	var ct = "application/json"
	//fmt.Println(bearer)
	bodd := bytes.NewBuffer(createbody())
	//fmt.Println(bodd)

	req1, err := http.NewRequest("PUT", url, bodd)
	req1.Header.Add("Authorization", bearer)
	req1.Header.Add("Content-Type", ct)

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req1)
	time.Sleep(10 * time.Second)

	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	rb, _ := ioutil.ReadAll(resp.Body)

	return err, resp.Status, rb
}

func Checkextstat(token string, subid string, rgname string, vmname string) (err error, st string, m map[string]string) {
	var extr Extension
	m = make(map[string]string)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	url := "https://management.azure.com/subscriptions/" + subid + "/resourceGroups/" + rgname + "/providers/Microsoft.Compute/virtualMachines/" + vmname + "?api-version=2018-10-01"

	var bearer = "Bearer " + token
	var ct = "application/json"
	//fmt.Println(bearer)

	//fmt.Println(bodd)

	req1, err := http.NewRequest("GET", url, nil)
	req1.Header.Add("Authorization", bearer)
	req1.Header.Add("Content-Type", ct)

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req1)

	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	rb, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(rb, &extr)

	for _, rt := range extr.Resources {
		m[rt.Properties.Type] = rt.Properties.ProvisioningState

	}

	return err, resp.Status, m
}

func Delrequest(token string, subid string, rgname string, vmname string) (err error, st string) {

	url := "https://management.azure.com/subscriptions/" + subid + "/resourceGroups/" + rgname + "/providers/Microsoft.Compute/virtualMachines/" + vmname + "/extensions/enablevmaccess?api-version=2018-10-01"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var bearer = "Bearer " + token
	var ct = "application/json"
	//fmt.Println(bearer)
	//bodd := bytes.NewBuffer(createbody())
	//fmt.Println(bodd)

	req1, err := http.NewRequest("DELETE", url, nil)
	req1.Header.Add("Authorization", bearer)
	req1.Header.Add("Content-Type", ct)

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req1)

	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	//rb, _ := ioutil.ReadAll(resp.Body)
	time.Sleep(20 * time.Second)

	return err, resp.Status
}
