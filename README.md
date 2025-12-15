# 💱 Currency Converter

Terminal tabanlı, güzel arayüzlü para birimi dönüştürücü. Bubble Tea (TUI) ve Cobra CLI ile geliştirilmiştir.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)

## ✨ Özellikler

- 🎨 **Güzel Terminal Arayüzü** - Bubble Tea ile interaktif TUI
- 🔄 **Gerçek Zamanlı Kurlar** - ExchangeRate-API ile güncel döviz kurları
- 💻 **Çift Mod** - İnteraktif TUI veya hızlı CLI komutları
- 🔍 **Para Birimi Filtreleme** - Yazarak arama yapabilme
- 🌍 **150+ Para Birimi** - Tüm dünya para birimleri desteklenir

## 🚀 Kurulum

### Gereksinimler

- Go 1.21 veya üzeri
- ExchangeRate-API API Key (ücretsiz)

### API Key Alma

1. [ExchangeRate-API](https://www.exchangerate-api.com/) sitesine gidin
2. Ücretsiz hesap oluşturun
3. Dashboard'dan API key'inizi kopyalayın

### Projeyi Kurma

```bash
# Repoyu klonlayın
git clone https://github.com/yourusername/Currency-Converter.git
cd Currency-Converter

# Bağımlılıkları yükleyin
go mod tidy

# .env dosyası oluşturun
cp .env.example .env

# .env dosyasını düzenleyin ve API key'inizi ekleyin
# EXCHANGE_RATE_API_KEY=your_actual_api_key_here
```

### Derleme

```bash
go build -o currency-converter .
```

## 📖 Kullanım

### İnteraktif Mod (TUI)

```bash
./currency-converter
```

Bu komut ile interaktif terminal arayüzü açılır:
- ↑/↓ tuşları ile para birimi seçin
- Klavyeden yazarak filtreleme yapın
- Enter ile seçin
- Miktar girin ve dönüşümü görün

### Hızlı Dönüşüm (CLI)

```bash
# 100 USD'yi EUR'ya çevir
./currency-converter -f USD -t EUR -a 100

# Diğer örnekler
./currency-converter --from TRY --to USD --amount 1000
./currency-converter -f GBP -t JPY -a 50
```

### Yardım

```bash
./currency-converter --help
```

## 🔧 Yapılandırma

### Environment Variables

| Değişken | Açıklama | Zorunlu |
|----------|----------|---------|
| `EXCHANGE_RATE_API_KEY` | ExchangeRate-API API Key | ✅ Evet |

### .env Dosyası

```bash
# .env.example dosyasını kopyalayın
cp .env.example .env

# API key'inizi ekleyin
EXCHANGE_RATE_API_KEY=your_api_key_here
```

> ⚠️ **Önemli:** `.env` dosyasını asla Git'e commit etmeyin! `.gitignore` zaten bunu engelleyecek şekilde ayarlanmıştır.

## 📁 Proje Yapısı

```
Currency-Converter/
├── cmd/
│   └── root.go           # Cobra CLI komutları
├── internal/
│   ├── api/
│   │   └── client.go     # ExchangeRate-API client
│   ├── config/
│   │   └── config.go     # Yapılandırma yönetimi
│   └── ui/
│       └── tui.go        # Bubble Tea TUI
├── .env.example          # API key şablonu
├── .gitignore            # Git ignore listesi
├── go.mod                # Go modül dosyası
├── main.go               # Giriş noktası
└── README.md             # Bu dosya
```

## 🛠️ Teknolojiler

- [Go](https://golang.org/) - Programlama dili
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [ExchangeRate-API](https://www.exchangerate-api.com/) - Döviz kuru API

## 📝 Lisans

MIT License - Detaylar için [LICENSE](LICENSE) dosyasına bakın.

## 🤝 Katkıda Bulunma

Pull request'ler kabul edilir. Büyük değişiklikler için önce bir issue açarak ne yapmak istediğinizi tartışalım.