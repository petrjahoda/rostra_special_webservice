function displayTable(TableData) {
    let orderCode = [];
    let orderName = [];
    let productName = [];
    let sytelineWorkplace = []
    let orderStart = []
    let orderRequestedTotal = []
    let totalProducedCount = []
    let terminalInputOrderProducedCount = []
    let totalTransferredCount = []
    let terminalInputOrderTransferredCount = []
    let waitingForTransferCount = []
    let totalNokCount = []
    for (data of TableData) {
        orderCode.push(data.OrderCode)
        orderName.push(data.OrderName)
        productName.push(data.ProductName)
        sytelineWorkplace.push(data.SytelineWorkplace)
        orderStart.push(data.OrderStart)
        orderRequestedTotal.push(data.OrderRequestedTotal)
        totalProducedCount.push(data.TotalProducedCount)
        terminalInputOrderProducedCount.push(data.TerminalInputOrderProducedCount)
        totalTransferredCount.push(data.TotalTransferredCount)
        terminalInputOrderTransferredCount.push(data.TerminalInputOrderTransferredCount)
        waitingForTransferCount.push(data.WaitingForTransferCount)
        totalNokCount.push(data.TotalNokCount)
    }
    let tableRef = document.getElementById('table').getElementsByTagName('tbody')[0];
    for (let index = 0; index < orderCode.length; index++) {
        tableRef.insertRow().innerHTML =
            "<td>" + orderCode[index] + "</td>" +
            "<td>" + orderName[index] + "</td>" +
            "<td>" + productName[index] + "</td>" +
            "<td>" + sytelineWorkplace[index] + "</td>" +
            "<td>" + orderStart[index] + "</td>" +
            "<td>" + orderRequestedTotal[index] + "</td>" +
            "<td>" + totalProducedCount[index] + "</td>" +
            "<td>" + terminalInputOrderProducedCount[index] + "</td>" +
            "<td>" + totalTransferredCount[index] + "</td>" +
            "<td>" + terminalInputOrderTransferredCount[index] + "</td>" +
            "<td>" + waitingForTransferCount[index] + "</td>" +
            "<td>" + totalNokCount[index] + "</td>";
    }
    for (let i = 0; i < table.rows.length; i++) {
        table.rows[i].addEventListener('click', function () {
            let row = "";
            for (let j = 0; j < this.cells.length; j++) {
                row += this.cells[j].innerHTML + ";";
            }
            processTableRow(row)
        });
    }
}

function processTableRow(row) {
    let dataSplitted = row.split(";")
    let orderFromTable = dataSplitted[1]
    let orderSplitted = orderFromTable.split("-")
    let orderInputData = orderSplitted[0]
    let operationInputData = orderSplitted[1]
    console.log("Order Input data: " + orderInputData)
    console.log("Operation Input data: " + operationInputData)
    let workplace = dataSplitted[3]
    let workplaceSplitted = workplace.split(";")
    let workplaceCode = workplaceSplitted[0]
    let workplaceName = workplaceSplitted[1]
    console.log("Workplace Input data: " + workplaceCode)
    let data = {OrderInput: orderInputData, UserInput: sessionStorage.getItem("userInput")};
    fetch("/check_order_input", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            if (result.Result === "ok") {
                orderInput.value = result.OrderName;
                orderRow.classList.add("disabled");
                operationRow.classList.remove("disabled");
                orderOkButton.disabled = true;
                orderBackButton.disabled = true;
                operationOkButton.disabled = false;
                operationBackButton.disabled = false;
                operationSelect.placeholder = "Zadejte číslo operace výrobního příkazu";
                sessionStorage.setItem("orderId", result.OrderId)
                infoOrderId.textContent = result.OrderId
                sessionStorage.setItem("orderName", result.OrderName)
                infoOrderName.textContent = result.OrderName
                sessionStorage.setItem("orderInput", result.OrderInput)
                infoOrderInput.textContent = result.OrderInput
                sessionStorage.setItem("priznakSeriovaVyroba", result.PriznakSeriovaVyroba)
                infoOrderPriznakSeriovaVyroba.textContent = result.PriznakSeriovaVyroba
                sessionStorage.setItem("productId", result.ProductId)
                let tableData = {};
                for (let operation of result.Operations) {
                    if (operationInputData === operation.Operace) {
                        tableData[operation.Operace] = operation.Operace + ": " + operation.Pracoviste + " [" + operation.PracovistePopis + "]"
                    }
                }
                const select = Metro.getPlugin("#operation-select", 'select');
                select.data({
                    "Načtené operace": tableData
                });
                infoRostra.textContent = ""
                data = {
                    OperationSelect: operationSelect.value,
                    OrderInput: sessionStorage.getItem("orderInput"),
                    ProductId: sessionStorage.getItem("productId"),
                    UserInput: sessionStorage.getItem("userInput")
                };
                fetch("/check_operation_input", {
                    method: "POST",
                    body: JSON.stringify(data)
                }).then((response) => {
                    response.text().then(function (data) {
                        let result = JSON.parse(data);
                        if (result.Result === "ok") {
                            operationOkButton.disabled = true
                            operationBackButton.disabled = true
                            workplaceOkButton.disabled = false
                            workplaceBackButton.disabled = false
                            operationRow.classList.add("disabled")
                            workplaceRow.classList.remove("disabled")
                            sessionStorage.setItem("operationValue", operationSelect.value)
                            infoOperationInput.textContent = operationSelect.value
                            sessionStorage.setItem("parovyDil", result.ParovyDil)
                            infoOperationParovyDil.textContent = result.ParovyDil
                            sessionStorage.setItem("seznamParovychDilu", result.SeznamParovychDilu)
                            infoOperationSeznamParovychDilu.textContent = result.SeznamParovychDilu
                            sessionStorage.setItem("jenPrenosMnozstvi", result.JenPrenosMnozstvi)
                            infoOperationJenPrenosMnozstvi.textContent = result.JenPrenosMnozstvi
                            sessionStorage.setItem("priznakMn2", result.PriznakMn2)
                            infoOperationPriznakMn2.textContent = result.PriznakMn2
                            sessionStorage.setItem("priznakMn3", result.PriznakMn3)
                            infoOperationPriznakMn3.textContent = result.PriznakMn3
                            sessionStorage.setItem("mn2Ks", result.Mn2Ks)
                            infoOperationMn2Ks.textContent = result.Mn2Ks
                            sessionStorage.setItem("mn3Ks", result.Mn3Ks)
                            infoOperationMn3Ks.textContent = result.Mn3Ks
                            sessionStorage.setItem("priznakNasobnost", result.PriznakNasobnost)
                            infoOperationPriznakNasobnost.textContent = result.PriznakNasobnost
                            sessionStorage.setItem("nasobnost", result.Nasobnost)
                            infoOperationNasobnost.textContent = result.Nasobnost
                            infoOrderId.textContent = result.OrderId
                            sessionStorage.setItem("orderId", result.OrderId)
                            let tableData = {};
                            savedWorkplaces = result.Workplaces;
                            for (let workplace of result.Workplaces) {
                                if (workplace.ZapsiZdroj.includes(workplaceCode)) {
                                    tableData[workplace.ZapsiZdroj] = workplace.ZapsiZdroj
                                }
                            }
                            const select = Metro.getPlugin("#workplace-select", 'select');
                            select.data({
                                "Načtené pracoviště": tableData
                            });
                            sessionStorage.setItem("orderId", result.OrderId)
                            infoRostra.textContent = ""


                            for (let workplace of savedWorkplaces) {
                                if (workplace.ZapsiZdroj.includes(workplaceCode)) {
                                    console.log("Workplace found: " + workplaceSelect.value);
                                    workplaceBackButton.disabled = true
                                    workplaceOkButton.disabled = true
                                    sessionStorage.setItem("workplaceCode", workplaceSelect.value.split(";")[0])
                                    infoWorkplaceCode.textContent = workplaceSelect.value.split(";")[0]
                                    infoWorkplaceName.textContent = workplaceSelect.value.split(";")[1]
                                    sessionStorage.setItem("typZdrojeZapsi", workplace.TypZdrojeZapsi)
                                    infoWorkplaceTypZdrojeZapsi.textContent = workplace.TypZdrojeZapsi
                                    sessionStorage.setItem("viceVp", workplace.ViceVp)
                                    infoWorkplaceViceVp.textContent = workplace.ViceVp
                                    sessionStorage.setItem("priznakMn1", workplace.PriznakMn1)
                                    infoWorkplacePriznakMn1.textContent = workplace.PriznakMn1
                                    let data = {
                                        WorkplaceCode: sessionStorage.getItem("workplaceCode"),
                                        UserId: sessionStorage.getItem("userId"),
                                        OrderInput: sessionStorage.getItem("orderInput"),
                                        OperationSelect: sessionStorage.getItem("operationValue"),
                                        ParovyDil: sessionStorage.getItem("parovyDil"),
                                        SeznamParovychDilu: sessionStorage.getItem("seznamParovychDilu"),
                                        JenPrenosMnozstvi: sessionStorage.getItem("jenPrenosMnozstvi"),
                                        TypZdrojeZapsi: sessionStorage.getItem("typZdrojeZapsi"),
                                        ViceVp: sessionStorage.getItem("viceVp"),
                                        UserInput: sessionStorage.getItem("userInput")
                                    };
                                    fetch("/check_workplace_input", {
                                        method: "POST",
                                        body: JSON.stringify(data)
                                    }).then((response) => {
                                        response.text().then(function (data) {
                                            let result = JSON.parse(data);
                                            if (result.Result === "ok") {
                                                workplaceRow.classList.add("disabled")
                                                if (result.OkInput === "true") {
                                                    okRow.classList.remove("disabled")
                                                    countBackButton.disabled = false
                                                    countButton.disabled = false
                                                    countOkInput.focus()
                                                }
                                                if (result.NokInput === "true") {
                                                    nokRow.classList.remove("disabled")
                                                    countBackButton.disabled = false
                                                    countButton.disabled = false
                                                    let tableData = {};
                                                    for (let nokType of result.NokTypes) {
                                                        tableData[nokType.Kod + ";" + nokType.Nazev.replaceAll(" ", "")] = nokType.Kod + ";" + nokType.Nazev.replaceAll(" ", "")
                                                    }
                                                    const select = Metro.getPlugin("#nok-type-select", 'select');
                                                    select.data({
                                                        "Načtené neshody": tableData
                                                    });
                                                }
                                                if (result.StartButton === "true") {
                                                    startOrderButton.disabled = false
                                                }
                                                if (result.EndButton === "true") {
                                                    endOrderButton.disabled = false
                                                }
                                                if (result.TransferButton === "true") {
                                                    transferOrderButton.disabled = false
                                                }
                                                if (result.ClovekSelection === "true") {
                                                    clovekRadio.disabled = false
                                                }
                                                if (result.SerizeniSelection === "true") {
                                                    serizeniRadio.disabled = false
                                                }
                                                if (result.StrojSelection === "true") {
                                                    strojRadio.disabled = false
                                                }
                                                infoRostra.textContent = ""
                                            } else {
                                                infoRostra.text = result.WorkplaceError
                                            }

                                        });
                                    }).catch((error) => {
                                        infoRostra.textContent = error.toString()
                                    });
                                }
                            }
                        } else {
                            infoRostra.textContent = result.OperationError;
                        }
                    });
                }).catch((error) => {
                    infoRostra.textContent = error.toString()
                });


            } else {
                infoRostra.textContent = result.OrderError
                orderInput.placeholder = result.OrderError
                orderInput.value = ""
            }
        });
    }).catch((error) => {
        infoRostra.textContent = error.toString()
    });
}



