package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"zzz/internal/config"
	"zzz/internal/dto"
)

type TaskRepository interface {
	StopActiveTasks(ctx context.Context) error
	CreateTask(ctx context.Context) (string, error)
	CompleteTask(ctx context.Context, taskID string, status string) error
}

type DataRepository interface {
	SaveProducts(ctx context.Context, products []dto.ProductSource) error
	SaveClients(ctx context.Context, clients []dto.ClientSource) error
}

type TaskService struct {
	taskRepo   TaskRepository
	dataRepo   DataRepository
	cfg        *config.Config
	client     *http.Client
	cancelFunc context.CancelFunc
	mu         sync.Mutex
}

func NewTaskService(taskRepo TaskRepository, dataRepo DataRepository, cfg *config.Config) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
		dataRepo: dataRepo,
		cfg:      cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *TaskService) StartSync(ctx context.Context) (string, error) {
	s.mu.Lock()

	if s.cancelFunc != nil {
		s.cancelFunc()
	}

	bgCtx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel
	s.mu.Unlock()

	if err := s.taskRepo.StopActiveTasks(ctx); err != nil {
		return "", fmt.Errorf("failed to stop active tasks: %w", err)
	}

	taskID, err := s.taskRepo.CreateTask(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create task: %w", err)
	}

	go s.processData(bgCtx, taskID)

	return taskID, nil
}

func (s *TaskService) processData(ctx context.Context, taskID string) {
	defer func() {
		s.mu.Lock()
		if s.cancelFunc != nil {
			s.cancelFunc = nil
		}
		s.mu.Unlock()
	}()

	var wg sync.WaitGroup
	productsUrls := []string{s.cfg.SOURCE1, s.cfg.SOURCE2, s.cfg.SOURCE3}
	clientUrl := s.cfg.SOURCE4

	doneCh := make(chan struct{})

	go func() {
		for _, u := range productsUrls {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				s.fetchAndSaveProducts(ctx, url)
			}(u)
		}
		wg.Wait()

		if ctx.Err() == nil {
			s.fetchAndSaveClients(ctx, clientUrl)
		}

		close(doneCh)
	}()

	status := "completed"
	select {
	case <-ctx.Done():
		status = "stopped"
		log.Printf("Task %s stopped (interrupted by new request)", taskID)
	case <-doneCh:
		log.Printf("Task %s completed successfully", taskID)
	}

	if err := s.taskRepo.CompleteTask(context.Background(), taskID, status); err != nil {
		log.Printf("failed to complete task %s: %v", taskID, err)
	}
}
func (s *TaskService) fetchAndSaveProducts(ctx context.Context, url string) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("failed to fetch products: %v", err)
		return
	}

	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("failed to do request %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code %d: %s", resp.Status, url)
		return
	}
	var products []dto.ProductSource
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		log.Printf("failed to decode products: %s: %v", url, err)
		return
	}

	if err := s.dataRepo.SaveProducts(ctx, products); err != nil {
		log.Printf("failed to save products: %s: %v", url, err)
	}
}

func (s *TaskService) fetchAndSaveClients(ctx context.Context, url string) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Failed to create request for clients %s: %v", url, err)
		return
	}

	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("Failed to execute request for clients %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code %d from %s", resp.StatusCode, url)
		return
	}

	var clients []dto.ClientSource
	if err := json.NewDecoder(resp.Body).Decode(&clients); err != nil {
		log.Printf("Failed to decode clients JSON from %s: %v", url, err)
		return
	}

	if err := s.dataRepo.SaveClients(ctx, clients); err != nil {
		log.Printf("Failed to save clients from %s to DB: %v", url, err)
	}
}
