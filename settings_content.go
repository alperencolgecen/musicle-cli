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
kaynaklardan indirip yerel olarak yönetmesine yarar.`
}
