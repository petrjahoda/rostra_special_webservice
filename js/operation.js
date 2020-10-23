operationBackButton.addEventListener("click", () => {
    orderInput.value = ""
    orderInput.placeholder = "Zadejte číslo výrobního příkazu";
    const select = Metro.getPlugin("#operation-select", 'select');
    select.data({});
    orderRow.classList.remove("disabled");
    operationRow.classList.add("disabled")
    operationBackButton.disabled = true
    operationOkButton.disabled = true;
    orderOkButton.disabled = false;
    orderBackButton.disabled = false;
    infoRostra.textContent = ""
    infoOperationNasobnost.textContent = ""
    infoOperationPriznakNasobnost.textContent = ""
    infoOperationMn2Ks.textContent = ""
    infoOperationPriznakMn2.textContent = ""
    infoOperationMn3Ks.textContent = ""
    infoOperationPriznakMn3.textContent = ""
    infoOperationJenPrenosMnozstvi.textContent = ""
    infoOperationSeznamParovychDilu.textContent = ""
    infoOperationParovyDil.textContent = ""
    infoOperationInput.textContent = ""
    orderInput.focus()
})

operationOkButton.addEventListener("click", () => {
    processOperationInput();
})

operationSelect.addEventListener("keyup", function (event) {
    if (event.code === "Enter") {
        processOperationInput();
    }
});

function processOperationInput() {
    console.log("Operation selected: " + operationSelect.value);
    let data = {
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
                    tableData[workplace.ZapsiZdroj] = workplace.ZapsiZdroj
                }
                const select = Metro.getPlugin("#workplace-select", 'select');
                select.data({
                    "Načtené pracoviště": tableData
                });
                infoRostra.textContent = ""
                workplaceSelect.focus()
            } else {
                infoRostra.textContent = result.OperationError;
            }
        });
    }).catch((error) => {
        infoRostra.textContent = error.toString()
    });
}