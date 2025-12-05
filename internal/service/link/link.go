package link

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/SenechkaP/link-checker/internal/model"
	linkModels "github.com/SenechkaP/link-checker/internal/model/link"
	"github.com/jung-kurt/gofpdf"
)

type LinkService struct {
	wg          *sync.WaitGroup
	mu          sync.RWMutex
	sendedLinks map[int]linkModels.LinksResponse
	nextNum     int
	client      *http.Client
}

func NewLinkService(wg *sync.WaitGroup) *LinkService {
	return &LinkService{
		wg:          wg,
		mu:          sync.RWMutex{},
		sendedLinks: make(map[int]linkModels.LinksResponse),
		nextNum:     1,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *LinkService) GetStatuses(ctx context.Context, req linkModels.LinksRequest) (*linkModels.LinksResponse, error) {
	if len(req.Links) == 0 {
		return nil, model.ErrNoLinks
	}

	resp := &linkModels.LinksResponse{
		Links: make(map[string]string),
	}

	for _, url := range req.Links {
		s.wg.Go(func() {
			select {
			case <-ctx.Done():
				return
			default:
			}

			status := s.getStatus(ctx, url)

			s.mu.Lock()
			resp.Links[url] = status
			s.mu.Unlock()
		})
	}

	s.wg.Wait()

	s.mu.Lock()
	resp.LinksNum = s.nextNum
	s.sendedLinks[resp.LinksNum] = *resp
	s.nextNum++
	s.mu.Unlock()

	return resp, nil
}

func (s *LinkService) getStatus(ctx context.Context, url string) string {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "not_available"
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "not_available"
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return "available"
	}

	return "not_available"
}

func (s *LinkService) GetStatusesByNums(ctx context.Context, req linkModels.LinksListRequest) ([]byte, error) {
	if len(req.LinksList) == 0 {
		return nil, model.ErrNoNums
	}

	collected := make([]linkModels.LinksResponse, 0, len(req.LinksList))

	s.mu.RLock()
	for _, num := range req.LinksList {
		if num < 0 {
			s.mu.RUnlock()
			return nil, model.ErrNumLessZero
		}
		links, ok := s.sendedLinks[num]
		if !ok {
			s.mu.RUnlock()
			return nil, fmt.Errorf("num %d is not found", num)
		}
		collected = append(collected, links)
	}
	s.mu.RUnlock()

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 10)
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	for i, resp := range collected {
		title := fmt.Sprintf("Links set #%d (LinksNum=%d)", i+1, resp.LinksNum)
		pdf.SetFont("Arial", "B", 12)
		pdf.CellFormat(0, 7, title, "", 1, "L", false, 0, "")
		pdf.SetFont("Arial", "", 11)

		for site, status := range resp.Links {
			line := fmt.Sprintf("%s: %s", site, status)
			pdf.MultiCell(200.6, 6, line, "", "L", false)
		}

		pdf.Ln(4)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
