package store

import (
 "crypto/rand"; "encoding/hex"; "encoding/json"; "fmt"; "net/url"; "os"; "sort"; "strings"; "sync"; "time"
 "github.com/nono-vos/broker/internal/model"
)
type snapshot struct { Targets map[string]model.Target `json:"targets"`; Leases map[string]model.Lease `json:"leases"`; Events []model.Event `json:"events"` }
type Store struct { mu sync.RWMutex; path string; targets map[string]model.Target; leases map[string]model.Lease; events []model.Event }
func New(path string) *Store { s:=&Store{path:path,targets:map[string]model.Target{},leases:map[string]model.Lease{}}; s.load(); return s }
func (s *Store) load(){ if s.path=="" {return}; b,e:=os.ReadFile(s.path); if e!=nil{return}; var x snapshot; if json.Unmarshal(b,&x)==nil { if x.Targets!=nil{s.targets=x.Targets}; if x.Leases!=nil{s.leases=x.Leases}; s.events=x.Events } }
func (s *Store) saveLocked(){ if s.path=="" {return}; b,_:=json.MarshalIndent(snapshot{s.targets,s.leases,s.events},"","  "); tmp:=s.path+".tmp"; if os.WriteFile(tmp,b,0600)==nil {_=os.Rename(tmp,s.path)} }
func (s *Store) AddTarget(t model.Target) error { if t.ID==""||t.Name==""||len(t.Origins)==0{return fmt.Errorf("id, name and origins are required")}; for _,o:=range t.Origins {u,e:=url.Parse(o); if e!=nil||u.Scheme!="http"&&u.Scheme!="https"||u.Host=="" {return fmt.Errorf("invalid origin: %s",o)}}; s.mu.Lock(); defer s.mu.Unlock(); s.targets[t.ID]=t; s.saveLocked(); return nil }
func (s *Store) Targets() []model.Target { s.mu.RLock(); defer s.mu.RUnlock(); out:=make([]model.Target,0,len(s.targets)); for _,t:=range s.targets{out=append(out,t)}; sort.Slice(out,func(i,j int)bool{return out[i].ID<out[j].ID}); return out }
func (s *Store) GetTarget(id string)(model.Target,bool){s.mu.RLock();defer s.mu.RUnlock();t,ok:=s.targets[id];return t,ok}
func (s *Store) NewLease(targetID string,ttl time.Duration,budget int)(model.Lease,error){if _,ok:=s.GetTarget(targetID);!ok{return model.Lease{},fmt.Errorf("unknown target")};if ttl<=0||ttl>24*time.Hour{return model.Lease{},fmt.Errorf("ttl out of range")};if budget<=0||budget>100000{return model.Lease{},fmt.Errorf("budget out of range")};b:=make([]byte,12);if _,e:=rand.Read(b);e!=nil{return model.Lease{},e};l:=model.Lease{ID:"lease_"+hex.EncodeToString(b),TargetID:targetID,ExpiresAt:time.Now().Add(ttl),Budget:budget};s.mu.Lock();s.leases[l.ID]=l;s.saveLocked();s.mu.Unlock();return l,nil}
func (s *Store) GetLease(id string)(model.Lease,bool){s.mu.RLock();defer s.mu.RUnlock();l,ok:=s.leases[id];return l,ok}
func (s *Store) RevokeLease(id string) error{s.mu.Lock();defer s.mu.Unlock();l,ok:=s.leases[id];if !ok{return fmt.Errorf("unknown lease")};l.Revoked=true;s.leases[id]=l;s.saveLocked();return nil}
func (s *Store) RecordEvent(e model.Event) error{s.mu.Lock();defer s.mu.Unlock();l,ok:=s.leases[e.LeaseID];if !ok{return fmt.Errorf("unknown lease")};if l.Revoked||time.Now().After(l.ExpiresAt){return fmt.Errorf("lease inactive")};if l.Used>=l.Budget{return fmt.Errorf("lease budget exceeded")};l.Used++;s.leases[e.LeaseID]=l;e.ID="evt_"+fmt.Sprint(time.Now().UnixNano());e.CreatedAt=time.Now().UTC();s.events=append(s.events,e);if len(s.events)>10000{s.events=s.events[len(s.events)-10000:]};s.saveLocked();return nil}
func(s *Store) Events(lease string)[]model.Event{s.mu.RLock();defer s.mu.RUnlock();out:=[]model.Event{};for _,e:=range s.events{if lease==""||e.LeaseID==lease{out=append(out,e)}};return out}
func OriginAllowed(t model.Target, raw string) bool { u,e:=url.Parse(raw);if e!=nil{return false}; for _,o:=range t.Origins {v,_:=url.Parse(o);if strings.EqualFold(v.Scheme,u.Scheme)&&strings.EqualFold(v.Host,u.Host){return true}};return false }
