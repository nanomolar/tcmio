package controllers

import (
	"fmt"
	"tcmio/models"

	"github.com/astaxie/beego/orm"
)

func (this *MainController) ListLigands() {
	var tars []models.Ligand
	var tar models.Ligand
	offset, _ := this.GetInt64("offset")
	limit, _ := this.GetInt64("limit")
	fmt.Println(offset)
	fmt.Println(limit)
	_, err := tar.Query().Limit(limit, offset).All(&tars)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tars)

	this.responseMsg.SuccessMsg("", tars)
	this.Data["json"] = this.responseMsg
	this.ServeJSON()

}

func (this *MainController) DetailLigand() {
	var tar models.Ligand
	id := this.Ctx.Input.Param(":id")
	fmt.Println(id)
	err := tar.Query().Filter("id", id).One(&tar)
	if err != nil {
		fmt.Println(err)
		this.responseMsg.ErrorMsg("Not find", "")
		this.Data["json"] = this.responseMsg
		this.ServeJSON()
	}
	fmt.Println(tar)

	this.responseMsg.SuccessMsg("", tar)
	this.Data["json"] = this.responseMsg
	this.ServeJSON()

}

// structure search ligands
func (this *MainController) SearchLigands() {
	method := this.Ctx.Input.Param(":method")
	fmt.Println(method)
	//limit, _ := this.GetInt64("limit")
	//offset, _ := this.GetInt64("offset")
	limit := this.GetString("limit")
	offset := this.GetString("offset")
	if limit == "" {
		limit = "10"
	}
	if offset == "" {
		offset = "0"
	}
	query := this.GetString("query")
	url := ""
	if method == "sim" {
		threshold := this.GetString("threshold")
		url = "select * from ligand where mol @ (" + threshold + ", 1.0, '" + query + "', 'Tanimoto')::bingo.sim;"
	} else if method == "sub" {
		url = "select *  from ligand where mol @ ('" + query + "','')::bingo.sub" + " limit " + limit + " offset " + offset + ";"
	} else if method == "exact" {
		url = "select *  from ligand where mol @ ('" + query + "','')::bingo.exact" + " limit " + limit + " offset " + offset + ";"
	} else {
		this.responseMsg.ErrorMsg("method not support", "")
		this.Data["json"] = this.responseMsg
		this.ServeJSON()
	}

	fmt.Println(url)
	o := orm.NewOrm()
	var hits []models.Ligand
	//_, err := o.Raw("select * from ligand where mol @ (0.5, 1.0, 'Fc1ccccc1Cn2ccc(NC(=O)c3ccc(COc4ccccc4Cl)cc3)n2', 'Tanimoto')::bingo.sim limit 10 offset 0;").QueryRows(&hits)

	//_, err := o.Raw("select * from ligand where mol @ (0.5, 1.0, 'C1=CC=CC=C1', 'Tanimoto')::bingo.sim limit 10 offset 0;").QueryRows(&hits)
	_, err := o.Raw(url).QueryRows(&hits)
	if err != nil {
		fmt.Println(err)
		this.responseMsg.ErrorMsg("Not find", "")
		this.Data["json"] = this.responseMsg
		this.ServeJSON()
	}
	fmt.Println(hits)

	this.responseMsg.SuccessMsg("", hits)
	this.Data["json"] = this.responseMsg
	this.ServeJSON()
}
