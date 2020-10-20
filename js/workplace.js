workplaceBackButton.addEventListener("click", () => {
    operationRow.classList.remove("disabled");
    workplaceRow.classList.add("disabled")
    const select = Metro.getPlugin("#workplace-select", 'select');
    select.data({});
    workplaceOkButton.disabled = true
    workplaceBackButton.disabled = true
    operationBackButton.disabled = false
    operationOkButton.disabled = false
    infoRostra.textContent = ""
    infoError.textContent = ""
    infoWorkplaceName.textContent = ""
    infoWorkplacePriznakMn1.textContent = ""
    infoWorkplaceViceVp.textContent = ""
    infoWorkplaceTypZdrojeZapsi.textContent = ""
    infoWorkplaceCode.textContent = ""
    operationSelect.focus()
})

workplaceOkButton.addEventListener("click", () => {
    processWorkplaceInput();
})


workplaceSelect.addEventListener("keyup", function (event) {
    if (event.code === "Enter") {
        processWorkplaceInput();
    }
});

function processWorkplaceInput() {
    console.log("Workplace selected: " + workplaceSelect.value);
    for (workplace of savedWorkplaces) {
        if (workplace.ZapsiZdroj === workplaceSelect.value) {
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
                ViceVp: sessionStorage.getItem("viceVp")
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
                            let chyby = {};
                            for (nokType of result.NokTypes) {
                                chyby[nokType.Nazev] = nokType.Kod + ";" + nokType.Nazev + ""
                            }
                            const select = Metro.getPlugin("#nok-type-select", 'select');
                            select.data({
                                "Načtené neshody": chyby
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
                    } else {
                        infoError.text = result.WorkplaceError
                    }
                });
            }).catch((error) => {
                infoError.textContent = error.toString()
            });
        } else {
            infoError.textContent = "Pracoviště nebylo nalezeno"
        }
    }

}
