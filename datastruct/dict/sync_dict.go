package dict

import (
	"sync"
)

type SyncDict struct {
  m sync.Map
}

func MakeSyncDict() *SyncDict {
  return &SyncDict{}
}

func (s *SyncDict) Get(keys string) (val interface{}, exists bool){
  val, ok := s.m.Load(keys)
  return val, ok
}

func (s *SyncDict) Len() int {
  length := 0
  s.m.Range(func(key, value interface{}) bool {
    length++
    return true
  })

  return length
}

func (s *SyncDict) Put(key string, val interface{}) (result int) {
  _, ok := s.m.Load(key)
  s.m.Store(key, val)
  if ok {
    return 0
  }
  return 1
}

func (s *SyncDict) PutIfAbsent(key string, val interface{})(result int) {
  _, ok := s.m.Load(key)
  if ok {
    return 0
  }

  s.m.Store(key, val)
  return 1
}

func (s *SyncDict) PutIfExists(key string, val interface{})(result int) {
  _, ok := s.m.Load(key)
  if ok {
    s.m.Store(key, val)
    return 1
  }

  return 0
}
  
func (s *SyncDict) Remove(key string)(result int){
  _, ok := s.m.Load(key)
  s.m.Delete(key)
  if ok {
    return 1
  }

  return 0
}
  
func (s *SyncDict) ForEach(consumer Consumer) {
  s.m.Range(func(key, value interface{}) bool {
    consumer(key.(string), value)
    return true
  })
}
  
func (s *SyncDict) Keys() []string {
  result := make([]string, s.Len())
  i := 0
  s.m.Range(func(key, val interface{}) bool{
    result[i] = key.(string)
    i++;
    return true
  })

  return result
}
  
func (s *SyncDict) RandomKeys(limit int) []string {
  result := make([]string, s.Len())
  for i := 0; i < limit; i++ {
    s.m.Range(func(key, val interface{}) bool{
      result[i] = key.(string)
      return false
    })
  }

  return result
}
  
func (s *SyncDict) RandomDistinctKeys(limit int) []string {
  result := make([]string, s.Len())
  i := 0  
  s.m.Range(func(key, val interface{}) bool{
    result[i] = key.(string)
    i++
    if i == limit {
      return false
    }

    return true
  })

  return result
}
  
func  (s *SyncDict) clear() {
  *s = *MakeSyncDict()
}
