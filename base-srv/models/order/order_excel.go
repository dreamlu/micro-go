package order

import (
	"demo/base-srv/models/admin/setup/logistic"
	json2 "encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/log"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/type/json"
	"github.com/dreamlu/gt/tool/type/time"
	"io"
	"strconv"
	time2 "time"
)

// 订单导出
type OrderExport struct {
	OrderGoods
	OrderOrderNum   string     `json:"order_order_num"`
	OrderPaytime    time.CTime `json:"order_paytime"`
	OrderClientName string     `json:"order_client_name" gt:"field:order.client_name"`
	OrderPhone      string     `json:"order_phone"`
	OrderAddress    string     `json:"order_address"`
	OrderStatus     int8       `json:"order_status"` // 订单状态
	GoodsCode       string     `json:"goods_code"`
	GoodsName       string     `json:"goods_name"`
	BuyClientName   string     `json:"buy_client_name" gt:"field:client.name"` // 买家账号
}

type OrderNorm struct {
	Type string `json:"type"`
}

// 订单导入
// {"com": "huitongkuaidi", "createtime": "2020-05-08 11:44", "courier_company": "韵达快递", "tracking_number": "557003195853545"}
type OrderShipP struct {
	Com            string `json:"com"`
	CourierCompany string `json:"courier_company"`
	TrackingNumber string `json:"tracking_number"`
	Createtime     string `json:"createtime"`
}

func (o OrderShipP) String() string {
	b, err := json2.Marshal(o)
	if err != nil {
		log.Error(b)
		return ""
	}
	return string(b)
}

type OrderShip struct {
	OrderShipP
	OrderNum string `json:"order_num"`
}

// 通过excel批量发货
func (c *Order) ShipExcel(r io.Reader) (err error) {

	f, err := excelize.OpenReader(r)
	if err != nil {
		println(err.Error())
		return err
	}

	var osps []*OrderShip
	rows, err := f.GetRows("Sheet1")
	for k, row := range rows {
		if k < 1 {
			continue
		}
		var osp OrderShip
		osp.OrderNum = row[0]
		osp.CourierCompany = row[1]
		//osp.Com = row[2]
		osp.TrackingNumber = row[2]
		osp.Createtime = time2.Now().Format("2006-01-02 15:04")

		// 根据快递公司查找编码,懒得合并减少查询(待优化)
		var lo logistic.Logistic
		err = lo.GetByName(osp.CourierCompany)
		if err != nil {
			log.Error("批量导入excel,查询快递公司编码问题:", err.Error())
			continue
		}
		osp.Com = lo.Com

		osps = append(osps, &osp)
		//for _, colCell := range row {
		//	//print(colCell, "\t")
		//	osp.
		//}
		//println()
	}

	// 批量修改待发货状态的
	var s = int8(2)
	for _, v := range osps {
		var or = Order{
			//OrderNum:     v.OrderNum,
			LogisticInfo: json.CJSON(v.OrderShipP.String()),
			Status:       &s,
		}
		newOr, er := or.GetByOrderNum(v.OrderNum)
		if er != nil {
			log.Error("excel批量发货修改问题", er.Error())
			continue
		}
		or.ID = newOr.ID
		_, er = or.Update(&or)
		if er != nil {
			log.Error("excel批量发货修改问题", er.Error())
		}
	}

	return
}

// 订单导出
func (c *Order) ExportExcel(params cmap.CMap) (f *excelize.File, err error) {

	statusWhereSQL := ""
	status := params.Get("status")
	if status != "" {
		switch status {
		case "4", "5", "6":
			statusWhereSQL = "(`order_goods`.status = " + status + ")"
		default:
			statusWhereSQL = "(`order`.status = " + status + " and `order_goods`.status not in (4,5,6))"
		}
		params.Del("status")
	}

	// 查询基础订单
	var datas []*OrderExport
	crud2 := gt.NewCrud(
		gt.Model(OrderExport{}),
		gt.Data(&datas),
		gt.InnerTable([]string{"order_goods", "order"}),
		gt.LeftTable([]string{"order_goods", "goods", "order", "client"}),
		gt.SubWhereSQL(statusWhereSQL),
	)
	cd := crud2.GetMoreBySearch(params)
	if cd.Error() != nil {
		//log.Log.Error(err.Error())
		err = cd.Error()
		return
	}

	// 单多规格判断
	for _, v := range datas {
		if v.GsNormID != 0 { // 多规格, 进行覆盖
			//var gs norm.GsNorm
			gt.NewCrud(gt.Data(&v)).Select("select code as goods_code from gs_norm where id = ?", v.GsNormID).Single()
		}
	}

	f = excelize.NewFile()
	// Create a new sheet.
	//index := f.NewSheet("Sheet2")
	//sheet := f.Sheet
	//// Set value of a cell.
	//_ = f.SetCellValue("Sheet2", "A2", "Hello world.")
	const St = "Sheet1"
	_ = f.SetCellValue(St, "A1", "订单编号")
	_ = f.SetCellValue(St, "B1", "商品状态")
	_ = f.SetCellValue(St, "C1", "付款时间")
	_ = f.SetCellValue(St, "D1", "买家账号")
	_ = f.SetCellValue(St, "E1", "收货人姓名")
	_ = f.SetCellValue(St, "F1", "手机号")
	_ = f.SetCellValue(St, "G1", "收货地址")
	_ = f.SetCellValue(St, "H1", "商品名称")
	_ = f.SetCellValue(St, "I1", "商品规格")
	_ = f.SetCellValue(St, "J1", "商品编码(SKU)")
	_ = f.SetCellValue(St, "K1", "数量")
	_ = f.SetCellValue(St, "L1", "商品单价") // 实付价格
	// 设置宽度样式
	_ = f.SetColWidth(St, "A", "B", 25)
	_ = f.SetColWidth(St, "B", "T", 15)
	//_ = f.SetColWidth(St, "B", "I", 18)
	// Set active sheet of the workbook.
	// f.SetActiveSheet(index)
	for k, v := range datas {
		var no OrderNorm
		//for k2, v2 := range v.OrderGoods {
		num := strconv.Itoa(k + 2)
		_ = f.SetCellValue(St, "A"+num, v.OrderOrderNum)
		switch v.Status {
		case 4, 5, 6:
		default:
			v.Status = v.OrderStatus
		}
		_ = f.SetCellValue(St, "B"+num, Status(v.Status)) // 订单状态
		_ = f.SetCellValue(St, "C"+num, v.OrderPaytime)
		_ = f.SetCellValue(St, "D"+num, v.BuyClientName)
		_ = f.SetCellValue(St, "E"+num, v.OrderClientName)
		_ = f.SetCellValue(St, "F"+num, v.OrderPhone)
		_ = f.SetCellValue(St, "G"+num, v.OrderAddress)
		_ = f.SetCellValue(St, "H"+num, v.GoodsName)
		_ = v.GsNorm.Struct(&no)
		_ = f.SetCellValue(St, "I"+num, no.Type)
		_ = f.SetCellValue(St, "J"+num, v.GoodsCode)
		_ = f.SetCellValue(St, "K"+num, v.Num)
		_ = f.SetCellValue(St, "L"+num, v.Money)
		//}
	}
	return
}
