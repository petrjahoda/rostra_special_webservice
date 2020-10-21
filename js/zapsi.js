const userRow = document.getElementById("user-row");
const orderRow = document.getElementById("order-row");
const operationRow = document.getElementById("operation-row");
const workplaceRow = document.getElementById("workplace-row");
const okRow = document.getElementById("ok-row");
const nokRow = document.getElementById("nok-row");

const resetButton = document.getElementById("reset-button");

const userInput = document.getElementById("user-input");
const orderInput = document.getElementById("order-input");
const countOkInput = document.getElementById("count-ok-input");
const countNokInput = document.getElementById("count-nok-input");
const operationSelect = document.getElementById("operation-select");
const workplaceSelect = document.getElementById("workplace-select");
const nokTypeSelect = document.getElementById("nok-type-select");

const userOkButton = document.getElementById("user-ok-button");
const orderOkButton = document.getElementById("order-ok-button");
const operationOkButton = document.getElementById("operation-ok-button");
const workplaceOkButton = document.getElementById("workplace-ok-button");
const countButton = document.getElementById("count-button");

const orderBackButton = document.getElementById("order-back-button");
const operationBackButton = document.getElementById("operation-back-button");
const workplaceBackButton = document.getElementById("workplace-back-button");
const countBackButton = document.getElementById("count-back-button");

const clovekRadio = document.getElementById("clovek-radio");
const serizeniRadio = document.getElementById("serizeni-radio");
const strojRadio = document.getElementById("stroj-radio");

const startOrderButton = document.getElementById("start-order-button");
const endOrderButton = document.getElementById("end-order-button");
const transferOrderButton = document.getElementById("transfer-order-button");

const table = document.getElementById("table");
const tableBody = document.getElementById("table-body");

let savedWorkplaces = {};

resetButton.addEventListener("click", () => {
    window.location.replace('');
})

sessionStorage.clear()


const infoUserInput = document.getElementById("info-user-input");
const infoUserName = document.getElementById("info-user-name");
const infoUserId = document.getElementById("info-user-id");

const infoOrderInput = document.getElementById("info-order-input");
const infoOrderName = document.getElementById("info-order-name");
const infoOrderId = document.getElementById("info-order-id");
const infoOrderPriznakSeriovaVyroba = document.getElementById("info-order-priznak-seriova-vyroba");


const infoOperationInput = document.getElementById("info-operation-input");
const infoOperationParovyDil = document.getElementById("info-operation-parovy-dil");
const infoOperationSeznamParovychDilu = document.getElementById("info-operation-seznam-parovych-dilu");
const infoOperationJenPrenosMnozstvi = document.getElementById("info-operation-jen-prenos-mnozstvi");
const infoOperationPriznakMn2 = document.getElementById("info-operation-priznak-mn2");
const infoOperationPriznakMn3 = document.getElementById("info-operation-priznak-mn3");
const infoOperationMn2Ks = document.getElementById("info-operation-mn2-ks");
const infoOperationMn3Ks = document.getElementById("info-operation-mn3-ks");
const infoOperationPriznakNasobnost = document.getElementById("info-operation-priznak-nasobnost");
const infoOperationNasobnost = document.getElementById("info-operation-nasobnost");

const infoWorkplaceCode = document.getElementById("info-workplace-code");
const infoWorkplaceName = document.getElementById("info-workplace-name");
const infoWorkplaceTypZdrojeZapsi = document.getElementById("info-workplace-typ-zdroje-zapsi");
const infoWorkplaceViceVp = document.getElementById("info-workplace-vice-vp");
const infoWorkplacePriznakMn1 = document.getElementById("info-workplace-priznak-mn1");


const infoError = document.getElementById("info-error");
const infoRostra = document.getElementById("info-rostra");

const userInputCell = document.getElementById("user-input-cell")
