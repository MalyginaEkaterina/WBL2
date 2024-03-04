package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

/*
Утилита wget

Реализовать утилиту wget с возможностью скачивать сайты целиком.
*/

func main() {
	maxDepth := flag.Int("max-depth", 10, "Maximum depth to download")
	targetDir := flag.String("dir", "output", "Directory to store files")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: <url>")
		os.Exit(1)
	}

	addr := flag.Args()[0]
	parsedAddr, err := url.Parse(addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't parse url %s: %v\n", addr, err)
		os.Exit(1)
	}
	baseURL := &url.URL{
		Scheme: parsedAddr.Scheme,
		Host:   parsedAddr.Host,
	}

	if err := os.MkdirAll(*targetDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Can't create directory %s: %v\n", *targetDir, err)
		os.Exit(1)
	}

	w := Wget{
		targetDir:  *targetDir,
		baseURL:    baseURL,
		downloaded: make(map[string]string),
	}

	if _, err := w.Download(addr, *maxDepth); err != nil {
		os.Exit(1)
	}
}

// Wget структура, содержащая целевую директорию, домен, мапу скачанных файлов и количество скачанных файлов
type Wget struct {
	// Целевая директория
	targetDir string
	// Домен
	baseURL *url.URL
	// Мап из адреса в имя скачанного файла
	downloaded map[string]string
	// Количество скачанных файлов
	count int
}

// Download скачивает документ по адресу addr и возвращает путь к локальному файлу, куда он сохранен, либо URL.
func (w *Wget) Download(addr string, depth int) (string, error) {
	if depth <= 0 {
		return addr, nil
	}
	if filename, ok := w.downloaded[addr]; ok {
		return filename, nil
	}
	parsedAddr, err := url.Parse(addr)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %v", addr, err)
	}
	if parsedAddr.Hostname() != w.baseURL.Hostname() {
		return addr, nil
	}

	resp, err := http.Get(addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Downloading %s failed: %v\n", addr, err)
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read %s failed: %v\n", addr, err)
		return "", err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("content-type")
	if contentType == "" {
		contentType = mime.TypeByExtension(path.Ext(parsedAddr.Path))
	}
	if contentType == "" {
		fmt.Fprintf(os.Stderr, "%s has no content type, skipping\n", addr)
		return "", nil
	}

	isHTML := false
	if strings.Contains(contentType, "text/html") {
		isHTML = true
	}

	outputName, err := w.getOutputNameFor(contentType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get output name %s failed: %v\n", addr, err)
		return "", err
	}
	w.downloaded[addr] = outputName

	outputPath := path.Join(w.targetDir, outputName)
	// если html, то рекурсивно обходим все ссылки, иначе просто скачиваем файл
	if isHTML {
		body, err = w.replaceLinks(body, depth-1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Replacing link in %s failed: %v\n", addr, err)
			return "", err
		}
	}
	if err := os.WriteFile(outputPath, body, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Writing %s to %s failed: %v\n", addr, outputName, err)
		return "", err
	}
	fmt.Fprintf(os.Stdout, "Wrote %s to %s\n", addr, outputName)

	return outputName, nil
}

// getOutputNameFor возвращает имя для файла по номеру и расширение согласно contentType
func (w *Wget) getOutputNameFor(contentType string) (string, error) {
	exts, err := mime.ExtensionsByType(contentType)
	if err != nil {
		return "", fmt.Errorf("failed to get extensions for content type %s: %v", contentType, err)
	}
	filename := fmt.Sprintf("%d%s", w.count, exts[0])
	w.count++
	return filename, nil
}

// replaceLinks проходит по html, для каждой ссылки вызывает Download и заменяет ссылку на локальную
func (w *Wget) replaceLinks(body []byte, depth int) ([]byte, error) {
	node, err := html.Parse(bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var findLinksAndReplace func(*html.Node)
	findLinksAndReplace = func(n *html.Node) {
		if n.Type == html.ElementNode {
			var attrName string
			switch n.DataAtom {
			case atom.A:
				attrName = "href"
			case atom.Link:
				attrName = "href"
			case atom.Img:
				attrName = "src"
			case atom.Script:
				attrName = "src"
			case atom.Embed:
				attrName = "src"
			case atom.Video:
				attrName = "src"
			case atom.Audio:
				attrName = "src"
			}
			if attrName != "" {
				for i, a := range n.Attr {
					if a.Key == attrName {
						// если это просто фрагмент, то пропускаем
						if !strings.HasPrefix(a.Val, "#") {
							link, fragment, err := w.getAbsLink(a.Val)
							// если какая-либо ошибка, просто пропускаем
							if err == nil {
								newLink, err := w.Download(link, depth)
								if err == nil {
									if fragment != "" {
										newLink = newLink + "#" + fragment
									}
									n.Attr[i].Val = newLink
								}
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findLinksAndReplace(c)
		}
	}
	findLinksAndReplace(node)
	buf := bytes.NewBuffer(nil)
	err = html.Render(buf, node)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Преобразование относительных ссылок в абсолютные
func (w *Wget) getAbsLink(l string) (string, string, error) {
	link, err := url.Parse(l)
	if err != nil {
		return "", "", err
	}
	if !link.IsAbs() {
		link = w.baseURL.ResolveReference(link)
	}
	fragment := link.Fragment
	link.Fragment = ""
	return link.String(), fragment, nil
}
