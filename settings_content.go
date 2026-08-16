package main

// policiesContent returns the Musicle usage policies shown in the
// Settings → Policies tab. It is intentionally a single string; the renderer
// word-wraps and makes it scrollable.
func policiesContent() string {
	return `Musicle Kullanım Politikaları

1. Kişisel kullanım
Musicle, yalnızca kişisel müzik arşivinizi oluşturmak ve yönetmek için
tasarlanmış ücretsiz, açık kaynaklı bir araçtır. Ticari amaçlı kullanım
önerilmez.

2. Telif hakkı ve yasal sorumluluk
İndirdiğiniz içeriklerin telif haklarına ve bulunduğunuz ülkenin yasalarına
uymak tamamen sizin sorumluluğunuzdadır. Musicle, içeriğin yasal durumunu
denetlemez ve herhangi bir yasal tavsiye vermez.

3. Dağıtım yasağı
İndirilen dosyaları başkalarına satmak, yeniden paylaşmak veya telif hakkı
sahibinin izni olmadan kamuya açmak yasaktır. Araç yalnızca kendi
kütüphaneniz içindir.

4. Hizmet sağlayıcıların şartları
Spotify ve YouTube gibi kaynakların kullanım koşullarına uymak zorundasınız.
Bu araç, ilgili platformların resmi API'lerine veya kamuya açık bağlantı
noktalarına bağımlıdır ve platformlar değişiklik yaptığında çalışma
davranışı değişebilir.

5. Gizlilik
Kimlik bilgileri (örn. tarayıcı çerezleri) yalnızca oturum boyunca yerel
olarak kullanılır; Musicle hiçbir kişisel veriyi harici sunucuya göndermez.
Yapılandırma yalnızca makinenizdeki config.json dosyasında saklanır.

6. Sorumluluk reddi
Bu yazılım "olduğu gibi" sunulur, garanti verilmez. Yanlış kullanımdan
kaynaklanan sonuçlardan geliştirici sorumlu tutulamaz.`
}

// aboutContent returns the project information shown in Settings → About.
// This base version covers founding/metadata; later edits extend it with
// purpose, features and third-party credits.
func aboutContent() string {
	return `Musicle Hakkında

Proje ne zaman başladı?
26 Mart 2026'da ilk kod yapısı oluşturularak geliştirmeye başlandı.

İlk commit ne zaman atıldı?
26 Mart 2026 — "File Structure Changes" (Alperen Çölgeçen).

Kim tarafından kurulup yönetiliyor?
Kurucu ve geliştirici: Alperen ÇÖLGEÇEN
GitHub: @colgecen

Bu araç neye hizmet eder?
Kullanıcının kendi kişisel müzik kütüphanesini Spotify ve YouTube gibi
kaynaklardan indirip yerel olarak yönetmesine yarar.

Neye hizmet eder (amaç)?
Mevcut müzik akış platformlarındaki parçaları kendi kontrolünüzdeki bir
yerel arşive dönüştürmek; çevrimdışı dinleme, düzenleme ve kalıcı saklama
sağlamak. Musicle bir akış servisi değildir, bir kişisel arşiv yöneticisidir.

Temel özellikler
- Spotify ve YouTube bağlantı/arama desteği ile parça ve çalma listesi
  indirme.
- Otomatik MP3 dönüştürme, ID3 etiketleme ve kapak resmi işleme.
- Profil ve çalma listesi yönetimi; şarkılar song_list.txt ile izlenir.
- Tema ve dil seçenekleri (çoklu dil desteği).
- Ses sekmesi: çıkış cihazı seçimi ve ses limiti ayarı.
- Bu politika ve hakkında ekranı dahil ayarlar yönetimi.
- Spektrum/görselleştirme ve kayan şarkı sözü gibi oynatıcı özellikleri.

Kullanılan 3. taraf uygulamalar ve kütüphaneler
Dış araçlar:
- yt-dlp: YouTube/Spotify bağlantılarından ses akışı çekmek için kullanılan
  açık kaynaklı indirme motoru.
- ffmpeg (statik): indirilen sesi MP3'e dönüştürme, yeniden örnekleme ve
  kapak/resim işleme için.
- Spotify: şarkı/çalma listesi meta verisi ve arama için veri kaynağı.

Go kütüphaneleri:
- charmbracelet/bubbletea + lipgloss: terminal kullanıcı arayüzü (TUI).
- gopxl/beep + ebitengine/oto: ses çalma ve ses efektleri.
- ncruces/zenity: masaüstü dosya/iletİşim kutuları.
- atotto/clipboard: pano desteği.
- dhowden/tag: ID3 etiket okuma/yazma.
- mewkiz/flac: FLAC çözümleme.
- Diğer bağımlılıklar ve lisanslar README.md ve LICENSE dosyalarında
  belgelenmiştir.`
}
