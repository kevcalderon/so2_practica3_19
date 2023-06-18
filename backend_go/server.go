package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

)

type RAM struct {
	TotalRam string `json:"totalram"`
	FreeRam  string `json:"freeram"`
}

type MemorySegment struct {
	StartAddress string `json:"start_address"`
	EndAddress   string `json:"end_address"`
}


func main() {

	http.HandleFunc("/api/ram", handleRequest)
	http.HandleFunc("/api/cpu", handleCPURequest)
	router.HandleFunc("/api/memoria/{folder}", getMemorySegments).Methods("GET")
	fmt.Println("Servidor en ejecución en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("sh", "-c", "cat /proc/ram_grupo19")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error al ejecutar el comando", http.StatusInternalServerError)
		return
	}

	outputRam := string(out[:])

	var ram RAM
	err = json.Unmarshal([]byte(outputRam), &ram)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error al analizar el JSON", http.StatusInternalServerError)
		return
	}

	// Establecer las cabeceras adecuadas para la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Devolver el JSON generado
	fmt.Fprintf(w, "%s", outputRam)
}

func handleCPURequest(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("sh", "-c", "cat /proc/cpu_grupo19")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error al ejecutar el comando", http.StatusInternalServerError)
		return
	}

	outputCpu := string(out[:])

	// Establecer las cabeceras adecuadas para la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Devolver el valor de outputCpu
	fmt.Fprintf(w, "%s", outputCpu)
}
func getMemorySegments(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	folder := params["folder"]

	filePath := fmt.Sprintf("/proc/%s/maps", folder)

	memorySegments, err := readMemorySegments(filePath)
	if err != nil {
		http.Error(w, "Error al leer los segmentos de memoria", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(memorySegments)
	if err != nil {
		http.Error(w, "Error al convertir los segmentos de memoria a JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func readMemorySegments(filePath string) ([]MemorySegment, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	memorySegments := make([]MemorySegment, 0)

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 1 {
			addressRange := strings.Split(fields[0], "-")
			if len(addressRange) == 2 {
				segment := MemorySegment{
					StartAddress: addressRange[0],
					EndAddress:   addressRange[1],
				}
				memorySegments = append(memorySegments, segment)
			}
		}
	}

	return memorySegments, nil
}