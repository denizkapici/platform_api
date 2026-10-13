# 🚀 Platform API - DevSecOps & GitOps Pipeline

Bu proje, Go dilinde geliştirilmiş bir REST API'nin modern DevSecOps ve GitOps prensipleri kullanılarak Oracle Cloud (ARM64) üzerindeki Kubernetes (k3s) kümesine tam otomatik dağıtımını (CI/CD) içeren uçtan uca bir altyapı mimarisidir.

## 🛠️ Kullanılan Teknolojiler (Tech Stack)

* **Uygulama:** Go (Golang)
* **Altyapı:** Oracle Cloud Infrastructure (OCI) - ARM64 (Ampere A1)
* **Konteyner Orkestrasyonu:** Kubernetes (k3s)
* **GitOps & CD:** ArgoCD
* **Sürekli Entegrasyon (CI):** GitHub Actions
* **Konteyner Kayıtçısı:** GitHub Container Registry (GHCR)
* **Güvenlik (Sec):** Gitleaks (Secret Scanning), Trivy (Vulnerability Scanning)
* **Gözlemlenebilirlik (Observability):** Prometheus, Grafana, Helm

---

## 🏗️ Mimari ve Çalışma Prensibi

Proje, **Single Source of Truth (Tek Gerçeklik Kaynağı)** olarak GitHub'ı kabul eden GitOps felsefesiyle tasarlanmıştır. Sunucuya manuel müdahale (SSH ile) tamamen ortadan kaldırılmıştır.

1. **Geliştirme & Push:** Geliştirici kodu yazar ve `main` dalına pushlar.
2. **Güvenlik Taraması (CI):** GitHub Actions tetiklenir. Kod içinde şifre/token sızıntısı var mı diye `Gitleaks`, bağımlılıklarda zafiyet var mı diye `Trivy` taraması yapılır.
3. **Multi-Arch Build (CI):** ARM64 mimarisine sahip Oracle sunucusu için kodlar QEMU/Buildx ile derlenir, Docker imajı oluşturulur ve `ghcr.io` (GitHub Packages) üzerine yüklenir.
4. **GitOps Senkronizasyonu (CD):** Kubernetes üzerinde çalışan ArgoCD, repodaki `k8s/` klasörünü sürekli dinler. Yeni bir imaj veya YAML yapılandırması gördüğünde bunu kümeye otomatik yansıtır (Sync).
5. **Gözlemlenebilirlik:** Helm ile kurulan Kube-Prometheus-Stack, uygulamanın bellek, CPU ve HTTP metriklerini sürekli kazır (scrape). Grafana üzerinden canlı izleme yapılır.

---

## 📂 Proje Yapısı

```text
platform_api/
├── .github/workflows/
│   └── ci.yml             # DevSecOps CI Pipeline (Gitleaks, Trivy, Docker Build&Push)
├── k8s/
│   └── deployment.yaml    # Kubernetes Deployment, LoadBalancer Service ve ServiceMonitor
├── Dockerfile             # Uygulama konteyner imaj şablonu
├── main.go                # Go REST API kaynak kodu
└── README.md
```

---

## 🌐 Açık Uç Noktalar (Exposed Services)

Altyapıdaki servisler Kubernetes LoadBalancer ve NodePort özellikleri ile dış dünyaya şu şekilde açılmıştır:

* **Platform API:** `http://<SERVER_IP>/` (Port 80)
* **ArgoCD UI:** `https://<SERVER_IP>/` (Port 443 - HTTPS)
* **Grafana Dashboards:** `http://<SERVER_IP>:30000/` (NodePort)

---

## 🚀 Yerel Geliştirme (Local Development)

Projeyi yerel bilgisayarınızda çalıştırmak için:

1. Repoyu klonlayın:
   ```bash
   git clone https://github.com/denizkapici/platform_api.git
   cd platform_api
   ```
2. Go uygulamasını çalıştırın:
   ```bash
   go run main.go
   ```
3. Veya Docker ile derleyin:
   ```bash
   docker build -t platform-api .
   docker run -p 8080:8080 platform-api
   ```

---
*Bu proje modern bulut mimarisi, uçtan uca otomasyon ve güvenlik entegrasyonu (Shift-Left Security) yeteneklerini sergilemek amacıyla geliştirilmiştir.*