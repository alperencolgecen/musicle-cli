# musicle-cli v1.3.0

Bu sürüm, uygulamaya tarayıcı tabanlı doğrudan bağlantı ("Connect")
özelliğini getirir ve bir dizi UI ile hata düzeltmesi içerir.

Önemli değişiklikler:

Connect (Bağlayıcı) özelliği:
  * Connect ekranı ana sayfaya taşındı; spektrum paleti ayarları eklendi
  * Platform seçim görünümü sabit marka renkleriyle eklendi
  * Bağlanma modalı 2s yükleyici ile eklendi
  * Tarama hatası artık tam ekran yerine küçük modal; Esc ile kapanır
  * Modal iyileştirmeleri ve kart logolarının tamamen kaldırılması
  * Spotify metadata artık yt-dlp ile çekiliyor (web kazıyıcı fallback)
  * Kart görselleri braille ile çiziliyor (her terminalde net görünüm)
  * Bağlayıcı ve ana ekran kartlarının yüksekliği olabildiğince küçültüldü

Durum ve profil:
  * State'e Source alanı ve ImportFromBrowser eklendi
  * Profil otomatik isimlendirme ve kaydetme
  * Playlist listesi onay butonlarıyla

UI:
  * Home F1 döngüsü: bar -> Spotify -> YouTube -> playlist -> songs
  * Build hataları giderildi (components import, renderLogo çakışması)

Diğer:
  * Spotify ve YouTube logo varlıkları eklendi
  * Bağımlılıklar, build scripti ve birim testleri güncellendi
