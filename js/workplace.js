workplaceBackButton.addEventListener("click", () => {
    operationRow.classList.remove("disabled");
    workplaceRow.classList.add("disabled")
    const select = Metro.getPlugin("#workplace-select", 'select');
    select.data({});
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
            sessionStorage.setItem("workplaceCode", workplaceSelect.value.split(";")[0])
            infoWorkplaceCode.textContent = workplaceSelect.value.split(";")[0]
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
                            countOkInput.focus()
                        }
                        if (result.NokInput === "true") {
                            nokRow.classList.remove("disabled")
                        }
                        if (result.StartButton === "true") {
                            startOrderButton.classList.remove("disabled")
                        }
                        if (result.EndButton === "true") {
                            endOrderButton.classList.remove("disabled")
                        }
                        if (result.TransferButton === "true") {
                            transferOrderButton.classList.remove("disabled")
                        }
                        if (result.ClovekSelection === "true") {
                            clovekRadio.classList.remove("disabled")
                        }
                        if (result.SerizeniSelection === "true") {
                            serizeniRadio.classList.remove("disabled")
                        }
                        if (result.StrojSelection === "true") {
                            strojRadio.classList.remove("disabled")
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
