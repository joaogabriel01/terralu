package main

import (
	"fmt"

	"github.com/joaogabriel01/terralu"
	"github.com/rivo/tview"
)

type TerraluProviderInfo struct {
	Alias     string `json:"alias" validate:"required"`
	Region    string `json:"region" validate:"required"`
	ApiKey    string `json:"api_key" validate:"required"`
	KeyID     string `json:"key_id"`
	KeySecret string `json:"key_secret"`
}

type AppData struct {
	terralu.TerraluProviderInfo
	Template string
}

type VMData struct {
	Name        string
	MachineType string
	Image       string
	SSHKeyName  string
}

var app *tview.Application
var pages *tview.Pages
var data *AppData = &AppData{}
var terraluProvider terralu.Terralu

func main() {

	app = tview.NewApplication()
	pages = tview.NewPages()
	fmt.Println("Terralu CLI")
	form := tview.NewForm().
		AddInputField("Api Key", "", 50, nil, func(text string) {
			data.ApiKey = text
		}).
		AddInputField("Key Id", "", 50, nil, func(text string) {
			data.KeyID = text
		}).
		AddInputField("Key Secret", "", 50, nil, func(text string) {
			data.KeySecret = text
		}).
		AddInputField("Region", "", 50, nil, func(text string) {
			data.Region = text
		}).
		AddInputField("Alias", "", 50, nil, func(text string) {
			data.Alias = text
		}).
		AddDropDown("Template", []string{"Nativo", "Customizado"}, 0, func(option string, optionIndex int) {
			data.Template = option
		}).
		AddButton("Save", func() {
			terraluProvider = terralu.NewTerralu(&data.TerraluProviderInfo)
			_, err := terraluProvider.GenerateTerraformGenericProviderConfig()
			if err != nil {
				panic("deu ruim")
			}

			chooseService()
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)

	pages.AddPage("main", form, true, true)

	if err := app.SetRoot(pages, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}

func chooseService() {
	form := tview.NewForm().
		AddButton("VMs", func() {
			vms()
		}).
		AddButton("MySQL", func() {
			showNotImplemented()
		}).
		AddButton("BlockStorage", func() {
			showNotImplemented()
		}).
		AddButton("ObjectStorage", func() {
			showNotImplemented()
		}).
		AddButton("Back", func() {
			pages.SwitchToPage("main")
		})

	form.SetBorder(true).SetTitle("Escolha o serviço").SetTitleAlign(tview.AlignLeft)

	pages.AddPage("chooseService", form, true, true)
	pages.SwitchToPage("chooseService")
}

func vms() {
	var vmData VMData

	form := tview.NewForm().
		AddInputField("Name", "", 50, nil, func(text string) {
			vmData.Name = text
		}).
		AddInputField("Machine Type", "", 50, nil, func(text string) {
			vmData.MachineType = text
		}).
		AddInputField("Image", "", 50, nil, func(text string) {
			vmData.Image = text
		}).
		AddInputField("SSH Key Name", "", 50, nil, func(text string) {
			vmData.SSHKeyName = text
		}).
		AddButton("Create", func() {
			showProvider(&vmData)
		}).
		AddButton("Back", func() {
			pages.SwitchToPage("chooseService")
		})

	form.SetBorder(true).SetTitle("Configure VM").SetTitleAlign(tview.AlignLeft)

	pages.AddPage("vms", form, true, true)
	pages.SwitchToPage("vms")
}

func showProvider(vmData *VMData) {
	required := terralu.VirtualMachineRequiredFields{
		Name:        vmData.Name,
		MachineType: &terralu.MachineTypeSchema{Name: vmData.MachineType},
		Image:       &terralu.ImageSchema{Name: vmData.Image},
		SSHKeyName:  vmData.SSHKeyName,
	}
	machine := terralu.VirtualMachineInstance{
		RequiredFields: required,
	}
	response, err := terraluProvider.GenerateTerraformVirtualMachineConfig(&machine)
	if err != nil {
		panic(err)
	}
	text := tview.NewTextView().
		SetText(response)

	text.SetBorder(true).SetTitle("VM Data").SetTitleAlign(tview.AlignLeft)
	terraluProvider.Save()
	pages.AddPage("vmData", text, true, true)
	pages.SwitchToPage("vmData")
}

func showNotImplemented() {
	modal := tview.NewModal().
		SetText("Esta funcionalidade não está implementada ainda.").
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.SwitchToPage("chooseService")
		})

	pages.AddPage("notImplemented", modal, true, true)
	pages.SwitchToPage("notImplemented")
}
