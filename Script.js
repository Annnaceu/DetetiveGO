// Função para abrir e fechar o menu hamburguer
function toggleMenu() {
    const menu = document.getElementById('menu');
    menu.classList.toggle('hidden');
}

// Função para mostrar/ocultar a seção de curiosidades
function toggleCuriosities() {
    const curiositiesContent = document.getElementById('curiositiesContent');
    curiositiesContent.classList.toggle('hidden');
}

// Submissão do formulário
window.onload = () => {
    document.getElementById("ipForm").addEventListener("submit", async function(event) {
        event.preventDefault();

        const ipBase = document.getElementById("ipBase").value;
        const deviceList = document.getElementById("deviceList");
        const errorDiv = document.getElementById("error");

        // Resetando resultados anteriores
        deviceList.innerHTML = '';
        errorDiv.innerHTML = '';

        try {
            const response = await fetch(`http://192.168.29.32:8080/scan?ip=${ipBase}`);

            if (!response.ok) {
                throw new Error("Nenhum dispositivo ativo encontrado ou erro na requisição");
            }

            const devices = await response.json();

            devices.forEach(device => {
                const li = document.createElement("li");
                li.textContent = `IP: ${device.ip}, Nome: ${device.name}`;
                deviceList.appendChild(li);
            });
        } catch (error) {
            errorDiv.textContent = "Erro ao buscar dispositivos: " + error.message;
            errorDiv.classList.remove('hidden');
        }
    });
};






