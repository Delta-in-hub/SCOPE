❯ paru -Ql ollama

/usr/bin/ollama
extern void llamaLog(int level, char* text, void* user_data);
uprobe:/usr/bin/ollama:llamaLog


uprobe:/usr/bin/ollama:llamaProgressCallback
uprobe:/usr/bin/ollama:sink


ollama /usr/lib/ollama/libggml-base.so

uprobe:/usr/lib/ollama/libggml-base.so:gguf_type_size(gguf_type)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_write_to_buf(gguf_context const*, std::vector<signed char, std::allocator<signed char> >&, bool)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_init_from_file_impl(_IO_FILE*, gguf_init_params)
uprobe:/usr/lib/ollama/libggml-base.so:bool gguf_read_emplace_helper<std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > >(gguf_reader const&, std::vector<gguf_kv, std::allocator<gguf_kv> >&, std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, bool, unsigned long)
uprobe:/usr/lib/ollama/libggml-base.so:bool gguf_read_emplace_helper<signed char>(gguf_reader const&, std::vector<gguf_kv, std::allocator<gguf_kv> >&, std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, bool, unsigned long)
uprobe:/usr/lib/ollama/libggml-base.so:bool gguf_read_emplace_helper<bool>(gguf_reader const&, std::vector<gguf_kv, std::allocator<gguf_kv> >&, std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, bool, unsigned long)
uprobe:/usr/lib/ollama/libggml-base.so:bool gguf_read_emplace_helper<double>(gguf_reader const&, std::vector<gguf_kv, std::allocator<gguf_kv> >&, std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, bool, unsigned long)
uprobe:/usr/lib/ollama/libggml-base.so:bool gguf_read_emplace_helper<float>(gguf_reader const&, std::vector<gguf_kv, std::allocator<gguf_kv> >&, std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, bool, unsigned long)
uprobe:/usr/lib/ollama/libggml-base.so:bool gguf_read_emplace_helper<unsigned char>(gguf_reader const&, std::vector<gguf_kv, std::allocator<gguf_kv> >&, std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, bool, unsigned long)
uprobe:/usr/lib/ollama/libggml-base.so:bool gguf_read_emplace_helper<int>(gguf_reader const&, std::vector<gguf_kv, std::allocator<gguf_kv> >&, std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, bool, unsigned long)
uprobe:/usr/lib/ollama/libggml-base.so:bool gguf_read_emplace_helper<unsigned int>(gguf_reader const&, std::vector<gguf_kv, std::allocator<gguf_kv> >&, std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, bool, unsigned long)
uprobe:/usr/lib/ollama/libggml-base.so:bool gguf_read_emplace_helper<long>(gguf_reader const&, std::vector<gguf_kv, std::allocator<gguf_kv> >&, std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, bool, unsigned long)
uprobe:/usr/lib/ollama/libggml-base.so:bool gguf_read_emplace_helper<unsigned long>(gguf_reader const&, std::vector<gguf_kv, std::allocator<gguf_kv> >&, std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, bool, unsigned long)
uprobe:/usr/lib/ollama/libggml-base.so:bool gguf_read_emplace_helper<short>(gguf_reader const&, std::vector<gguf_kv, std::allocator<gguf_kv> >&, std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, bool, unsigned long)
uprobe:/usr/lib/ollama/libggml-base.so:bool gguf_read_emplace_helper<unsigned short>(gguf_reader const&, std::vector<gguf_kv, std::allocator<gguf_kv> >&, std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, bool, unsigned long)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv(gguf_kv&&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >, std::allocator<std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > > > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<signed char>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<signed char, std::allocator<signed char> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<signed char>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, signed char)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<bool>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<bool, std::allocator<bool> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<bool>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, bool)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<double>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<double, std::allocator<double> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<double>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, double)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<float>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<float, std::allocator<float> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<float>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, float)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<unsigned char>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, unsigned char)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<int>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<int, std::allocator<int> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<int>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, int)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<unsigned int>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<unsigned int, std::allocator<unsigned int> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<unsigned int>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, unsigned int)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<long>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<long, std::allocator<long> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<long>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, long)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<unsigned long>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<unsigned long, std::allocator<unsigned long> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<unsigned long>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, unsigned long)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<short>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<short, std::allocator<short> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<short>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, short)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<unsigned short>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<unsigned short, std::allocator<unsigned short> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<unsigned short>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, unsigned short)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv(gguf_kv&&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >, std::allocator<std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > > > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<signed char>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<signed char, std::allocator<signed char> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<signed char>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, signed char)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<bool>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<bool, std::allocator<bool> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<bool>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, bool)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<double>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<double, std::allocator<double> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<double>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, double)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<float>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<float, std::allocator<float> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<float>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, float)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<unsigned char>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, unsigned char)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<int>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<int, std::allocator<int> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<int>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, int)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<unsigned int>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<unsigned int, std::allocator<unsigned int> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<unsigned int>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, unsigned int)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<long>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<long, std::allocator<long> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<long>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, long)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<unsigned long>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<unsigned long, std::allocator<unsigned long> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<unsigned long>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, unsigned long)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<short>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<short, std::allocator<short> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<short>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, short)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<unsigned short>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, std::vector<unsigned short, std::allocator<unsigned short> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::gguf_kv<unsigned short>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, unsigned short)
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::~gguf_kv()
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::~gguf_kv()
uprobe:/usr/lib/ollama/libggml-base.so:gguf_reader::read(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >&) const
uprobe:/usr/lib/ollama/libggml-base.so:gguf_kv::get_ne() const
uprobe:/usr/lib/ollama/libggml-base.so:std::mersenne_twister_engine<unsigned long, 32ul, 624ul, 397ul, 31ul, 2567483615ul, 11ul, 4294967295ul, 7ul, 2636928640ul, 15ul, 4022730752ul, 18ul, 1812433253ul>::_M_gen_rand()
uprobe:/usr/lib/ollama/libggml-base.so:std::map<gguf_type, char const*, std::less<gguf_type>, std::allocator<std::pair<gguf_type const, char const*> > >::~map()
uprobe:/usr/lib/ollama/libggml-base.so:std::map<gguf_type, char const*, std::less<gguf_type>, std::allocator<std::pair<gguf_type const, char const*> > >::~map()
uprobe:/usr/lib/ollama/libggml-base.so:std::map<gguf_type, unsigned long, std::less<gguf_type>, std::allocator<std::pair<gguf_type const, unsigned long> > >::~map()
uprobe:/usr/lib/ollama/libggml-base.so:std::map<gguf_type, unsigned long, std::less<gguf_type>, std::allocator<std::pair<gguf_type const, unsigned long> > >::~map()
uprobe:/usr/lib/ollama/libggml-base.so:void std::vector<gguf_tensor_info, std::allocator<gguf_tensor_info> >::_M_realloc_append<gguf_tensor_info const&>(gguf_tensor_info const&)
uprobe:/usr/lib/ollama/libggml-base.so:void std::vector<gguf_kv, std::allocator<gguf_kv> >::_M_realloc_append<char const*&, std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > >(char const*&, std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >&&)
uprobe:/usr/lib/ollama/libggml-base.so:void std::vector<gguf_kv, std::allocator<gguf_kv> >::_M_realloc_append<char const*&, std::vector<std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >, std::allocator<std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > > >&>(char const*&, std::vector<std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >, std::allocator<std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > > >&)
uprobe:/usr/lib/ollama/libggml-base.so:void std::vector<gguf_kv, std::allocator<gguf_kv> >::_M_realloc_append<char const*&, std::vector<signed char, std::allocator<signed char> >&>(char const*&, std::vector<signed char, std::allocator<signed char> >&)
uprobe:/usr/lib/ollama/libggml-base.so:void std::vector<gguf_kv, std::allocator<gguf_kv> >::_M_realloc_append<char const*&, signed char&>(char const*&, signed char&)
uprobe:/usr/lib/ollama/libggml-base.so:void std::vector<gguf_kv, std::allocator<gguf_kv> >::_M_realloc_append<char const*&, bool&>(char const*&, bool&)
uprobe:/usr/lib/ollama/libggml-base.so:void std::vector<gguf_kv, std::allocator<gguf_kv> >::_M_realloc_append<char const*&, double&>(char const*&, double&)
uprobe:/usr/lib/ollama/libggml-base.so:void std::vector<gguf_kv, std::allocator<gguf_kv> >::_M_realloc_append<char const*&, float&>(char const*&, float&)
uprobe:/usr/lib/ollama/libggml-base.so:void std::vector<gguf_kv, std::allocator<gguf_kv> >::_M_realloc_append<char const*&, unsigned char&>(char const*&, unsigned char&)
uprobe:/usr/lib/ollama/libggml-base.so:void std::vector<gguf_kv, std::allocator<gguf_kv> >::_M_realloc_append<char const*&, int&>(char const*&, int&)
uprobe:/usr/lib/ollama/libggml-base.so:void std::vector<gguf_kv, std::allocator<gguf_kv> >::_M_realloc_append<char const*&, unsigned int&>(char const*&, unsigned int&)
uprobe:/usr/lib/ollama/libggml-base.so:void std::vector<gguf_kv, std::allocator<gguf_kv> >::_M_realloc_append<char const*&, long&>(char const*&, long&)
uprobe:/usr/lib/ollama/libggml-base.so:void std::vector<gguf_kv, std::allocator<gguf_kv> >::_M_realloc_append<char const*&, unsigned long&>(char const*&, unsigned long&)
uprobe:/usr/lib/ollama/libggml-base.so:void std::vector<gguf_kv, std::allocator<gguf_kv> >::_M_realloc_append<char const*&, short&>(char const*&, short&)
uprobe:/usr/lib/ollama/libggml-base.so:void std::vector<gguf_kv, std::allocator<gguf_kv> >::_M_realloc_append<char const*&, unsigned short&>(char const*&, unsigned short&)
uprobe:/usr/lib/ollama/libggml-base.so:void std::vector<std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >, std::allocator<std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > > >::_M_realloc_append<std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&>(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:std::vector<std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >, std::allocator<std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > > >::~vector()
uprobe:/usr/lib/ollama/libggml-base.so:std::vector<std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >, std::allocator<std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > > >::~vector()
uprobe:/usr/lib/ollama/libggml-base.so:std::vector<signed char, std::allocator<signed char> >::_M_default_append(unsigned long)
uprobe:/usr/lib/ollama/libggml-base.so:void std::vector<signed char, std::allocator<signed char> >::_M_realloc_append<signed char const&>(signed char const&)
uprobe:/usr/lib/ollama/libggml-base.so:std::vector<long, std::allocator<long> >::_M_default_append(unsigned long)
uprobe:/usr/lib/ollama/libggml-base.so:std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >::_M_dispose()
uprobe:/usr/lib/ollama/libggml-base.so:std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >::_M_replace_cold(char*, unsigned long, char const*, unsigned long, unsigned long)
uprobe:/usr/lib/ollama/libggml-base.so:std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >::resize(unsigned long, char)
uprobe:/usr/lib/ollama/libggml-base.so:std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >::_M_assign(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&)
uprobe:/usr/lib/ollama/libggml-base.so:std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >::_M_mutate(unsigned long, unsigned long, char const*, unsigned long)
uprobe:/usr/lib/ollama/libggml-base.so:void std::shuffle<__gnu_cxx::__normal_iterator<long*, std::vector<long, std::allocator<long> > >, std::mersenne_twister_engine<unsigned long, 32ul, 624ul, 397ul, 31ul, 2567483615ul, 11ul, 4294967295ul, 7ul, 2636928640ul, 15ul, 4022730752ul, 18ul, 1812433253ul>&>(__gnu_cxx::__normal_iterator<long*, std::vector<long, std::allocator<long> > >, __gnu_cxx::__normal_iterator<long*, std::vector<long, std::allocator<long> > >, std::mersenne_twister_engine<unsigned long, 32ul, 624ul, 397ul, 31ul, 2567483615ul, 11ul, 4294967295ul, 7ul, 2636928640ul, 15ul, 4022730752ul, 18ul, 1812433253ul>&)
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_iq1_m
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_iq1_s
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_iq2_s
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_iq2_xs
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_iq2_xxs
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_iq3_s
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_iq3_xxs
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_iq4_nl
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_iq4_xs
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_q2_K
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_q3_K
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_q4_0
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_q4_1
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_q4_K
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_q5_0
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_q5_1
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_q5_K
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_q6_K
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_q8_0
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_q8_K
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_tq1_0
uprobe:/usr/lib/ollama/libggml-base.so:dequantize_row_tq2_0
uprobe:/usr/lib/ollama/libggml-base.so:ggml_abort
uprobe:/usr/lib/ollama/libggml-base.so:ggml_abs
uprobe:/usr/lib/ollama/libggml-base.so:ggml_abs_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_acc
uprobe:/usr/lib/ollama/libggml-base.so:ggml_acc_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_add
uprobe:/usr/lib/ollama/libggml-base.so:ggml_add1
uprobe:/usr/lib/ollama/libggml-base.so:ggml_add1_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_add_cast
uprobe:/usr/lib/ollama/libggml-base.so:ggml_add_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_add_rel_pos
uprobe:/usr/lib/ollama/libggml-base.so:ggml_add_rel_pos_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_aligned_free
uprobe:/usr/lib/ollama/libggml-base.so:ggml_aligned_malloc
uprobe:/usr/lib/ollama/libggml-base.so:ggml_arange
uprobe:/usr/lib/ollama/libggml-base.so:ggml_are_same_shape
uprobe:/usr/lib/ollama/libggml-base.so:ggml_are_same_stride
uprobe:/usr/lib/ollama/libggml-base.so:ggml_argmax
uprobe:/usr/lib/ollama/libggml-base.so:ggml_argsort
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_alloc_buffer
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_alloc_ctx_tensors
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_alloc_ctx_tensors_from_buft
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buffer_clear
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buffer_copy_tensor
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buffer_free
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buffer_get_alignment
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buffer_get_alloc_size
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buffer_get_base
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buffer_get_max_size
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buffer_get_size
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buffer_get_type
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buffer_get_usage
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buffer_init
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buffer_init_tensor
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buffer_is_host
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buffer_is_multi_buffer
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buffer_name
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buffer_reset
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buffer_set_usage
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buft_alloc_buffer
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buft_get_alignment
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buft_get_alloc_size
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buft_get_device
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buft_get_max_size
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buft_is_host
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_buft_name
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_compare_graph_backend
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_cpu_buffer_from_ptr
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_cpu_buffer_type
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_dev_backend_reg
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_dev_buffer_from_host_ptr
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_dev_buffer_type
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_dev_description
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_dev_get_props
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_dev_host_buffer_type
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_dev_init
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_dev_memory
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_dev_name
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_dev_offload_op
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_dev_supports_buft
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_dev_supports_op
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_dev_type
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_event_free
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_event_new
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_event_record
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_event_synchronize
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_event_wait
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_free
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_get_alignment
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_get_default_buffer_type
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_get_device
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_get_max_size
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_graph_compute
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_graph_compute_async
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_graph_copy
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_graph_copy_free
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_graph_plan_compute
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_graph_plan_create
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_graph_plan_free
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_guid
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_multi_buffer_alloc_buffer
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_multi_buffer_set_usage
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_name
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_offload_op
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_reg_dev_count
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_reg_dev_get
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_reg_get_proc_address
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_reg_name
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_sched_alloc_graph
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_sched_free
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_sched_get_backend
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_sched_get_buffer_size
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_sched_get_n_backends
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_sched_get_n_copies
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_sched_get_n_splits
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_sched_get_tensor_backend
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_sched_graph_compute
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_sched_graph_compute_async
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_sched_new
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_sched_reserve
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_sched_reset
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_sched_set_eval_callback
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_sched_set_tensor_backend
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_sched_synchronize
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_supports_buft
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_supports_op
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_synchronize
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_tensor_alloc
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_tensor_copy
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_tensor_copy_async
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_tensor_get
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_tensor_get_async
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_tensor_memset
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_tensor_set
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_tensor_set_async
uprobe:/usr/lib/ollama/libggml-base.so:ggml_backend_view_init
uprobe:/usr/lib/ollama/libggml-base.so:ggml_bf16_to_fp32
uprobe:/usr/lib/ollama/libggml-base.so:ggml_bf16_to_fp32_row
uprobe:/usr/lib/ollama/libggml-base.so:ggml_blck_size
uprobe:/usr/lib/ollama/libggml-base.so:ggml_build_backward_expand
uprobe:/usr/lib/ollama/libggml-base.so:ggml_build_forward_expand
uprobe:/usr/lib/ollama/libggml-base.so:ggml_can_repeat
uprobe:/usr/lib/ollama/libggml-base.so:ggml_cast
uprobe:/usr/lib/ollama/libggml-base.so:ggml_clamp
uprobe:/usr/lib/ollama/libggml-base.so:ggml_concat
uprobe:/usr/lib/ollama/libggml-base.so:ggml_cont
uprobe:/usr/lib/ollama/libggml-base.so:ggml_cont_1d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_cont_2d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_cont_3d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_cont_4d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_conv_1d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_conv_1d_dw
uprobe:/usr/lib/ollama/libggml-base.so:ggml_conv_1d_dw_ph
uprobe:/usr/lib/ollama/libggml-base.so:ggml_conv_1d_ph
uprobe:/usr/lib/ollama/libggml-base.so:ggml_conv_2d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_conv_2d_dw
uprobe:/usr/lib/ollama/libggml-base.so:ggml_conv_2d_s1_ph
uprobe:/usr/lib/ollama/libggml-base.so:ggml_conv_2d_sk_p0
uprobe:/usr/lib/ollama/libggml-base.so:ggml_conv_transpose_1d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_conv_transpose_2d_p0
uprobe:/usr/lib/ollama/libggml-base.so:ggml_cos
uprobe:/usr/lib/ollama/libggml-base.so:ggml_cos_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_count_equal
uprobe:/usr/lib/ollama/libggml-base.so:ggml_cpy
uprobe:/usr/lib/ollama/libggml-base.so:ggml_critical_section_end
uprobe:/usr/lib/ollama/libggml-base.so:ggml_critical_section_start
uprobe:/usr/lib/ollama/libggml-base.so:ggml_cross_entropy_loss
uprobe:/usr/lib/ollama/libggml-base.so:ggml_cross_entropy_loss_back
uprobe:/usr/lib/ollama/libggml-base.so:ggml_cycles
uprobe:/usr/lib/ollama/libggml-base.so:ggml_cycles_per_ms
uprobe:/usr/lib/ollama/libggml-base.so:ggml_diag
uprobe:/usr/lib/ollama/libggml-base.so:ggml_diag_mask_inf
uprobe:/usr/lib/ollama/libggml-base.so:ggml_diag_mask_inf_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_diag_mask_zero
uprobe:/usr/lib/ollama/libggml-base.so:ggml_diag_mask_zero_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_div
uprobe:/usr/lib/ollama/libggml-base.so:ggml_div_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_dup
uprobe:/usr/lib/ollama/libggml-base.so:ggml_dup_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_dup_tensor
uprobe:/usr/lib/ollama/libggml-base.so:ggml_element_size
uprobe:/usr/lib/ollama/libggml-base.so:ggml_elu
uprobe:/usr/lib/ollama/libggml-base.so:ggml_elu_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_exp
uprobe:/usr/lib/ollama/libggml-base.so:ggml_exp_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_flash_attn_back
uprobe:/usr/lib/ollama/libggml-base.so:ggml_flash_attn_ext
uprobe:/usr/lib/ollama/libggml-base.so:ggml_flash_attn_ext_get_prec
uprobe:/usr/lib/ollama/libggml-base.so:ggml_flash_attn_ext_set_prec
uprobe:/usr/lib/ollama/libggml-base.so:ggml_fopen
uprobe:/usr/lib/ollama/libggml-base.so:ggml_format_name
uprobe:/usr/lib/ollama/libggml-base.so:ggml_fp16_to_fp32
uprobe:/usr/lib/ollama/libggml-base.so:ggml_fp16_to_fp32_row
uprobe:/usr/lib/ollama/libggml-base.so:ggml_fp32_to_bf16
uprobe:/usr/lib/ollama/libggml-base.so:ggml_fp32_to_bf16_row
uprobe:/usr/lib/ollama/libggml-base.so:ggml_fp32_to_bf16_row_ref
uprobe:/usr/lib/ollama/libggml-base.so:ggml_fp32_to_fp16
uprobe:/usr/lib/ollama/libggml-base.so:ggml_fp32_to_fp16_row
uprobe:/usr/lib/ollama/libggml-base.so:ggml_free
uprobe:/usr/lib/ollama/libggml-base.so:ggml_ftype_to_ggml_type
uprobe:/usr/lib/ollama/libggml-base.so:ggml_gallocr_alloc_graph
uprobe:/usr/lib/ollama/libggml-base.so:ggml_gallocr_free
uprobe:/usr/lib/ollama/libggml-base.so:ggml_gallocr_get_buffer_size
uprobe:/usr/lib/ollama/libggml-base.so:ggml_gallocr_new
uprobe:/usr/lib/ollama/libggml-base.so:ggml_gallocr_new_n
uprobe:/usr/lib/ollama/libggml-base.so:ggml_gallocr_reserve
uprobe:/usr/lib/ollama/libggml-base.so:ggml_gallocr_reserve_n
uprobe:/usr/lib/ollama/libggml-base.so:ggml_gated_linear_attn
uprobe:/usr/lib/ollama/libggml-base.so:ggml_gelu
uprobe:/usr/lib/ollama/libggml-base.so:ggml_gelu_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_gelu_quick
uprobe:/usr/lib/ollama/libggml-base.so:ggml_gelu_quick_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_get_data
uprobe:/usr/lib/ollama/libggml-base.so:ggml_get_data_f32
uprobe:/usr/lib/ollama/libggml-base.so:ggml_get_first_tensor
uprobe:/usr/lib/ollama/libggml-base.so:ggml_get_max_tensor_size
uprobe:/usr/lib/ollama/libggml-base.so:ggml_get_mem_buffer
uprobe:/usr/lib/ollama/libggml-base.so:ggml_get_mem_size
uprobe:/usr/lib/ollama/libggml-base.so:ggml_get_name
uprobe:/usr/lib/ollama/libggml-base.so:ggml_get_next_tensor
uprobe:/usr/lib/ollama/libggml-base.so:ggml_get_no_alloc
uprobe:/usr/lib/ollama/libggml-base.so:ggml_get_rel_pos
uprobe:/usr/lib/ollama/libggml-base.so:ggml_get_rows
uprobe:/usr/lib/ollama/libggml-base.so:ggml_get_rows_back
uprobe:/usr/lib/ollama/libggml-base.so:ggml_get_tensor
uprobe:/usr/lib/ollama/libggml-base.so:ggml_get_type_traits
uprobe:/usr/lib/ollama/libggml-base.so:ggml_get_unary_op
uprobe:/usr/lib/ollama/libggml-base.so:ggml_graph_add_node
uprobe:/usr/lib/ollama/libggml-base.so:ggml_graph_clear
uprobe:/usr/lib/ollama/libggml-base.so:ggml_graph_cpy
uprobe:/usr/lib/ollama/libggml-base.so:ggml_graph_dump_dot
uprobe:/usr/lib/ollama/libggml-base.so:ggml_graph_dup
uprobe:/usr/lib/ollama/libggml-base.so:ggml_graph_get_grad
uprobe:/usr/lib/ollama/libggml-base.so:ggml_graph_get_grad_acc
uprobe:/usr/lib/ollama/libggml-base.so:ggml_graph_get_tensor
uprobe:/usr/lib/ollama/libggml-base.so:ggml_graph_n_nodes
uprobe:/usr/lib/ollama/libggml-base.so:ggml_graph_node
uprobe:/usr/lib/ollama/libggml-base.so:ggml_graph_nodes
uprobe:/usr/lib/ollama/libggml-base.so:ggml_graph_overhead
uprobe:/usr/lib/ollama/libggml-base.so:ggml_graph_overhead_custom
uprobe:/usr/lib/ollama/libggml-base.so:ggml_graph_print
uprobe:/usr/lib/ollama/libggml-base.so:ggml_graph_reset
uprobe:/usr/lib/ollama/libggml-base.so:ggml_graph_size
uprobe:/usr/lib/ollama/libggml-base.so:ggml_graph_view
uprobe:/usr/lib/ollama/libggml-base.so:ggml_group_norm
uprobe:/usr/lib/ollama/libggml-base.so:ggml_group_norm_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_guid_matches
uprobe:/usr/lib/ollama/libggml-base.so:ggml_hardsigmoid
uprobe:/usr/lib/ollama/libggml-base.so:ggml_hardswish
uprobe:/usr/lib/ollama/libggml-base.so:ggml_hash_set_free
uprobe:/usr/lib/ollama/libggml-base.so:ggml_hash_set_new
uprobe:/usr/lib/ollama/libggml-base.so:ggml_hash_set_reset
uprobe:/usr/lib/ollama/libggml-base.so:ggml_hash_size
uprobe:/usr/lib/ollama/libggml-base.so:ggml_im2col
uprobe:/usr/lib/ollama/libggml-base.so:ggml_im2col_back
uprobe:/usr/lib/ollama/libggml-base.so:ggml_init
uprobe:/usr/lib/ollama/libggml-base.so:ggml_is_3d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_is_contiguous
uprobe:/usr/lib/ollama/libggml-base.so:ggml_is_contiguous_0
uprobe:/usr/lib/ollama/libggml-base.so:ggml_is_contiguous_1
uprobe:/usr/lib/ollama/libggml-base.so:ggml_is_contiguous_2
uprobe:/usr/lib/ollama/libggml-base.so:ggml_is_empty
uprobe:/usr/lib/ollama/libggml-base.so:ggml_is_matrix
uprobe:/usr/lib/ollama/libggml-base.so:ggml_is_permuted
uprobe:/usr/lib/ollama/libggml-base.so:ggml_is_quantized
uprobe:/usr/lib/ollama/libggml-base.so:ggml_is_scalar
uprobe:/usr/lib/ollama/libggml-base.so:ggml_is_transposed
uprobe:/usr/lib/ollama/libggml-base.so:ggml_is_vector
uprobe:/usr/lib/ollama/libggml-base.so:ggml_leaky_relu
uprobe:/usr/lib/ollama/libggml-base.so:ggml_log
uprobe:/usr/lib/ollama/libggml-base.so:ggml_log_callback_default
uprobe:/usr/lib/ollama/libggml-base.so:ggml_log_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_log_internal
uprobe:/usr/lib/ollama/libggml-base.so:ggml_log_set
uprobe:/usr/lib/ollama/libggml-base.so:ggml_map_binary_f32
uprobe:/usr/lib/ollama/libggml-base.so:ggml_map_binary_inplace_f32
uprobe:/usr/lib/ollama/libggml-base.so:ggml_map_custom1
uprobe:/usr/lib/ollama/libggml-base.so:ggml_map_custom1_f32
uprobe:/usr/lib/ollama/libggml-base.so:ggml_map_custom1_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_map_custom1_inplace_f32
uprobe:/usr/lib/ollama/libggml-base.so:ggml_map_custom2
uprobe:/usr/lib/ollama/libggml-base.so:ggml_map_custom2_f32
uprobe:/usr/lib/ollama/libggml-base.so:ggml_map_custom2_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_map_custom2_inplace_f32
uprobe:/usr/lib/ollama/libggml-base.so:ggml_map_custom3
uprobe:/usr/lib/ollama/libggml-base.so:ggml_map_custom3_f32
uprobe:/usr/lib/ollama/libggml-base.so:ggml_map_custom3_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_map_custom3_inplace_f32
uprobe:/usr/lib/ollama/libggml-base.so:ggml_map_unary_f32
uprobe:/usr/lib/ollama/libggml-base.so:ggml_map_unary_inplace_f32
uprobe:/usr/lib/ollama/libggml-base.so:ggml_mean
uprobe:/usr/lib/ollama/libggml-base.so:ggml_mul
uprobe:/usr/lib/ollama/libggml-base.so:ggml_mul_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_mul_mat
uprobe:/usr/lib/ollama/libggml-base.so:ggml_mul_mat_id
uprobe:/usr/lib/ollama/libggml-base.so:ggml_mul_mat_set_prec
uprobe:/usr/lib/ollama/libggml-base.so:ggml_n_dims
uprobe:/usr/lib/ollama/libggml-base.so:ggml_nbytes
uprobe:/usr/lib/ollama/libggml-base.so:ggml_nbytes_pad
uprobe:/usr/lib/ollama/libggml-base.so:ggml_neg
uprobe:/usr/lib/ollama/libggml-base.so:ggml_neg_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_nelements
uprobe:/usr/lib/ollama/libggml-base.so:ggml_new_buffer
uprobe:/usr/lib/ollama/libggml-base.so:ggml_new_graph
uprobe:/usr/lib/ollama/libggml-base.so:ggml_new_graph_custom
uprobe:/usr/lib/ollama/libggml-base.so:ggml_new_tensor
uprobe:/usr/lib/ollama/libggml-base.so:ggml_new_tensor_1d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_new_tensor_2d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_new_tensor_3d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_new_tensor_4d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_norm
uprobe:/usr/lib/ollama/libggml-base.so:ggml_norm_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_nrows
uprobe:/usr/lib/ollama/libggml-base.so:ggml_op_desc
uprobe:/usr/lib/ollama/libggml-base.so:ggml_op_name
uprobe:/usr/lib/ollama/libggml-base.so:ggml_op_symbol
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_dataset_data
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_dataset_free
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_dataset_get_batch
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_dataset_init
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_dataset_labels
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_dataset_shuffle
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_default_params
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_epoch
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_epoch_callback_progress_bar
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_fit
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_forward
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_forward_backward
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_free
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_get_default_optimizer_params
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_grad_acc
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_init
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_inputs
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_labels
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_loss
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_ncorrect
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_outputs
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_pred
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_reset
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_result_accuracy
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_result_free
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_result_init
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_result_loss
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_result_ndata
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_result_pred
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_result_reset
uprobe:/usr/lib/ollama/libggml-base.so:ggml_opt_step_adamw
uprobe:/usr/lib/ollama/libggml-base.so:ggml_out_prod
uprobe:/usr/lib/ollama/libggml-base.so:ggml_pad
uprobe:/usr/lib/ollama/libggml-base.so:ggml_pad_reflect_1d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_permute
uprobe:/usr/lib/ollama/libggml-base.so:ggml_pool_1d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_pool_2d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_pool_2d_back
uprobe:/usr/lib/ollama/libggml-base.so:ggml_print_object
uprobe:/usr/lib/ollama/libggml-base.so:ggml_print_objects
uprobe:/usr/lib/ollama/libggml-base.so:ggml_quantize_chunk
uprobe:/usr/lib/ollama/libggml-base.so:ggml_quantize_free
uprobe:/usr/lib/ollama/libggml-base.so:ggml_quantize_init
uprobe:/usr/lib/ollama/libggml-base.so:ggml_quantize_requires_imatrix
uprobe:/usr/lib/ollama/libggml-base.so:ggml_relu
uprobe:/usr/lib/ollama/libggml-base.so:ggml_relu_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_repeat
uprobe:/usr/lib/ollama/libggml-base.so:ggml_repeat_back
uprobe:/usr/lib/ollama/libggml-base.so:ggml_reset
uprobe:/usr/lib/ollama/libggml-base.so:ggml_reshape
uprobe:/usr/lib/ollama/libggml-base.so:ggml_reshape_1d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_reshape_2d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_reshape_3d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_reshape_4d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_rms_norm
uprobe:/usr/lib/ollama/libggml-base.so:ggml_rms_norm_back
uprobe:/usr/lib/ollama/libggml-base.so:ggml_rms_norm_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_rope
uprobe:/usr/lib/ollama/libggml-base.so:ggml_rope_custom
uprobe:/usr/lib/ollama/libggml-base.so:ggml_rope_custom_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_rope_ext
uprobe:/usr/lib/ollama/libggml-base.so:ggml_rope_ext_back
uprobe:/usr/lib/ollama/libggml-base.so:ggml_rope_ext_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_rope_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_rope_multi
uprobe:/usr/lib/ollama/libggml-base.so:ggml_rope_multi_back
uprobe:/usr/lib/ollama/libggml-base.so:ggml_rope_yarn_corr_dims
uprobe:/usr/lib/ollama/libggml-base.so:ggml_row_size
uprobe:/usr/lib/ollama/libggml-base.so:ggml_rwkv_wkv6
uprobe:/usr/lib/ollama/libggml-base.so:ggml_scale
uprobe:/usr/lib/ollama/libggml-base.so:ggml_scale_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_set
uprobe:/usr/lib/ollama/libggml-base.so:ggml_set_1d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_set_1d_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_set_2d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_set_2d_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_set_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_set_input
uprobe:/usr/lib/ollama/libggml-base.so:ggml_set_loss
uprobe:/usr/lib/ollama/libggml-base.so:ggml_set_name
uprobe:/usr/lib/ollama/libggml-base.so:ggml_set_no_alloc
uprobe:/usr/lib/ollama/libggml-base.so:ggml_set_output
uprobe:/usr/lib/ollama/libggml-base.so:ggml_set_param
uprobe:/usr/lib/ollama/libggml-base.so:ggml_set_zero
uprobe:/usr/lib/ollama/libggml-base.so:ggml_sgn
uprobe:/usr/lib/ollama/libggml-base.so:ggml_sgn_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_sigmoid
uprobe:/usr/lib/ollama/libggml-base.so:ggml_sigmoid_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_silu
uprobe:/usr/lib/ollama/libggml-base.so:ggml_silu_back
uprobe:/usr/lib/ollama/libggml-base.so:ggml_silu_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_sin
uprobe:/usr/lib/ollama/libggml-base.so:ggml_sin_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_soft_max
uprobe:/usr/lib/ollama/libggml-base.so:ggml_soft_max_ext
uprobe:/usr/lib/ollama/libggml-base.so:ggml_soft_max_ext_back
uprobe:/usr/lib/ollama/libggml-base.so:ggml_soft_max_ext_back_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_soft_max_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_sqr
uprobe:/usr/lib/ollama/libggml-base.so:ggml_sqr_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_sqrt
uprobe:/usr/lib/ollama/libggml-base.so:ggml_sqrt_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_ssm_conv
uprobe:/usr/lib/ollama/libggml-base.so:ggml_ssm_scan
uprobe:/usr/lib/ollama/libggml-base.so:ggml_status_to_string
uprobe:/usr/lib/ollama/libggml-base.so:ggml_step
uprobe:/usr/lib/ollama/libggml-base.so:ggml_step_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_sub
uprobe:/usr/lib/ollama/libggml-base.so:ggml_sub_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_sum
uprobe:/usr/lib/ollama/libggml-base.so:ggml_sum_rows
uprobe:/usr/lib/ollama/libggml-base.so:ggml_tallocr_alloc
uprobe:/usr/lib/ollama/libggml-base.so:ggml_tallocr_new
uprobe:/usr/lib/ollama/libggml-base.so:ggml_tanh
uprobe:/usr/lib/ollama/libggml-base.so:ggml_tanh_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_tensor_overhead
uprobe:/usr/lib/ollama/libggml-base.so:ggml_threadpool_params_default
uprobe:/usr/lib/ollama/libggml-base.so:ggml_threadpool_params_init
uprobe:/usr/lib/ollama/libggml-base.so:ggml_threadpool_params_match
uprobe:/usr/lib/ollama/libggml-base.so:ggml_time_init
uprobe:/usr/lib/ollama/libggml-base.so:ggml_time_ms
uprobe:/usr/lib/ollama/libggml-base.so:ggml_time_us
uprobe:/usr/lib/ollama/libggml-base.so:ggml_timestep_embedding
uprobe:/usr/lib/ollama/libggml-base.so:ggml_top_k
uprobe:/usr/lib/ollama/libggml-base.so:ggml_transpose
uprobe:/usr/lib/ollama/libggml-base.so:ggml_type_name
uprobe:/usr/lib/ollama/libggml-base.so:ggml_type_size
uprobe:/usr/lib/ollama/libggml-base.so:ggml_type_sizef
uprobe:/usr/lib/ollama/libggml-base.so:ggml_unary
uprobe:/usr/lib/ollama/libggml-base.so:ggml_unary_inplace
uprobe:/usr/lib/ollama/libggml-base.so:ggml_unary_op_name
uprobe:/usr/lib/ollama/libggml-base.so:ggml_unpad
uprobe:/usr/lib/ollama/libggml-base.so:ggml_unravel_index
uprobe:/usr/lib/ollama/libggml-base.so:ggml_upscale
uprobe:/usr/lib/ollama/libggml-base.so:ggml_upscale_ext
uprobe:/usr/lib/ollama/libggml-base.so:ggml_used_mem
uprobe:/usr/lib/ollama/libggml-base.so:ggml_validate_row_data
uprobe:/usr/lib/ollama/libggml-base.so:ggml_view_1d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_view_2d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_view_3d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_view_4d
uprobe:/usr/lib/ollama/libggml-base.so:ggml_view_tensor
uprobe:/usr/lib/ollama/libggml-base.so:ggml_win_part
uprobe:/usr/lib/ollama/libggml-base.so:ggml_win_unpart
uprobe:/usr/lib/ollama/libggml-base.so:gguf_add_tensor
uprobe:/usr/lib/ollama/libggml-base.so:gguf_find_key
uprobe:/usr/lib/ollama/libggml-base.so:gguf_find_tensor
uprobe:/usr/lib/ollama/libggml-base.so:gguf_free
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_alignment
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_arr_data
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_arr_data_n
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_arr_n
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_arr_str
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_arr_type
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_data_offset
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_key
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_kv_type
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_meta_data
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_meta_size
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_n_kv
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_n_tensors
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_tensor_name
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_tensor_offset
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_tensor_size
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_tensor_type
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_val_bool
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_val_data
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_val_f32
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_val_f64
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_val_i16
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_val_i32
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_val_i64
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_val_i8
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_val_str
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_val_u16
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_val_u32
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_val_u64
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_val_u8
uprobe:/usr/lib/ollama/libggml-base.so:gguf_get_version
uprobe:/usr/lib/ollama/libggml-base.so:gguf_init_empty
uprobe:/usr/lib/ollama/libggml-base.so:gguf_init_from_file
uprobe:/usr/lib/ollama/libggml-base.so:gguf_remove_key
uprobe:/usr/lib/ollama/libggml-base.so:gguf_set_arr_data
uprobe:/usr/lib/ollama/libggml-base.so:gguf_set_arr_str
uprobe:/usr/lib/ollama/libggml-base.so:gguf_set_kv
uprobe:/usr/lib/ollama/libggml-base.so:gguf_set_tensor_data
uprobe:/usr/lib/ollama/libggml-base.so:gguf_set_tensor_type
uprobe:/usr/lib/ollama/libggml-base.so:gguf_set_val_bool
uprobe:/usr/lib/ollama/libggml-base.so:gguf_set_val_f32
uprobe:/usr/lib/ollama/libggml-base.so:gguf_set_val_f64
uprobe:/usr/lib/ollama/libggml-base.so:gguf_set_val_i16
uprobe:/usr/lib/ollama/libggml-base.so:gguf_set_val_i32
uprobe:/usr/lib/ollama/libggml-base.so:gguf_set_val_i64
uprobe:/usr/lib/ollama/libggml-base.so:gguf_set_val_i8
uprobe:/usr/lib/ollama/libggml-base.so:gguf_set_val_str
uprobe:/usr/lib/ollama/libggml-base.so:gguf_set_val_u16
uprobe:/usr/lib/ollama/libggml-base.so:gguf_set_val_u32
uprobe:/usr/lib/ollama/libggml-base.so:gguf_set_val_u64
uprobe:/usr/lib/ollama/libggml-base.so:gguf_set_val_u8
uprobe:/usr/lib/ollama/libggml-base.so:gguf_type_name
uprobe:/usr/lib/ollama/libggml-base.so:gguf_write_to_file
uprobe:/usr/lib/ollama/libggml-base.so:iq2xs_free_impl
uprobe:/usr/lib/ollama/libggml-base.so:iq2xs_init_impl
uprobe:/usr/lib/ollama/libggml-base.so:iq3xs_free_impl
uprobe:/usr/lib/ollama/libggml-base.so:iq3xs_init_impl
uprobe:/usr/lib/ollama/libggml-base.so:quantize_iq1_m
uprobe:/usr/lib/ollama/libggml-base.so:quantize_iq1_s
uprobe:/usr/lib/ollama/libggml-base.so:quantize_iq2_s
uprobe:/usr/lib/ollama/libggml-base.so:quantize_iq2_xs
uprobe:/usr/lib/ollama/libggml-base.so:quantize_iq2_xxs
uprobe:/usr/lib/ollama/libggml-base.so:quantize_iq3_s
uprobe:/usr/lib/ollama/libggml-base.so:quantize_iq3_xxs
uprobe:/usr/lib/ollama/libggml-base.so:quantize_iq4_nl
uprobe:/usr/lib/ollama/libggml-base.so:quantize_iq4_xs
uprobe:/usr/lib/ollama/libggml-base.so:quantize_q2_K
uprobe:/usr/lib/ollama/libggml-base.so:quantize_q3_K
uprobe:/usr/lib/ollama/libggml-base.so:quantize_q4_0
uprobe:/usr/lib/ollama/libggml-base.so:quantize_q4_1
uprobe:/usr/lib/ollama/libggml-base.so:quantize_q4_K
uprobe:/usr/lib/ollama/libggml-base.so:quantize_q5_0
uprobe:/usr/lib/ollama/libggml-base.so:quantize_q5_1
uprobe:/usr/lib/ollama/libggml-base.so:quantize_q5_K
uprobe:/usr/lib/ollama/libggml-base.so:quantize_q6_K
uprobe:/usr/lib/ollama/libggml-base.so:quantize_q8_0
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_iq2_s_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_iq3_s_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_iq3_xxs_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_iq4_nl_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_iq4_xs_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_q2_K_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_q3_K_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_q4_0_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_q4_1_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_q4_K_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_q5_0_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_q5_1_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_q5_K_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_q6_K_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_q8_0_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_q8_1_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_q8_K_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_tq1_0_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_row_tq2_0_ref
uprobe:/usr/lib/ollama/libggml-base.so:quantize_tq1_0
uprobe:/usr/lib/ollama/libggml-base.so:quantize_tq2_0



ollama /usr/lib/ollama/libggml-cpu-alderlake.so
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_backend_cpu_aarch64_buffer_type()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_backend_cpu_get_extra_buffers_type()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_backend_cpu_device_context::ggml_backend_cpu_device_context()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_backend_cpu_device_context::ggml_backend_cpu_device_context()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_backend_cpu_device_context::~ggml_backend_cpu_device_context()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_backend_cpu_device_context::~ggml_backend_cpu_device_context()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::tensor_traits::~tensor_traits()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::tensor_traits::~tensor_traits()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::tensor_traits::~tensor_traits()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::extra_buffer_type::~extra_buffer_type()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::extra_buffer_type::~extra_buffer_type()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::extra_buffer_type::~extra_buffer_type()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_q4_0, 4l, 4l>::compute_forward(ggml_compute_params*, ggml_tensor*)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_q4_0, 4l, 4l>::repack(ggml_tensor*, void const*, unsigned long)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_q4_0, 4l, 4l>::work_size(int, ggml_tensor const*, unsigned long&)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_q4_0, 4l, 4l>::~tensor_traits()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_q4_0, 4l, 4l>::~tensor_traits()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_q4_0, 4l, 4l>::~tensor_traits()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_q4_0, 8l, 4l>::compute_forward(ggml_compute_params*, ggml_tensor*)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_q4_0, 8l, 4l>::repack(ggml_tensor*, void const*, unsigned long)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_q4_0, 8l, 4l>::work_size(int, ggml_tensor const*, unsigned long&)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_q4_0, 8l, 4l>::~tensor_traits()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_q4_0, 8l, 4l>::~tensor_traits()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_q4_0, 8l, 4l>::~tensor_traits()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_q4_0, 8l, 8l>::compute_forward(ggml_compute_params*, ggml_tensor*)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_q4_0, 8l, 8l>::repack(ggml_tensor*, void const*, unsigned long)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_q4_0, 8l, 8l>::work_size(int, ggml_tensor const*, unsigned long&)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_q4_0, 8l, 8l>::~tensor_traits()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_q4_0, 8l, 8l>::~tensor_traits()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_q4_0, 8l, 8l>::~tensor_traits()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_iq4_nl, 4l, 4l>::compute_forward(ggml_compute_params*, ggml_tensor*)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_iq4_nl, 4l, 4l>::repack(ggml_tensor*, void const*, unsigned long)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_iq4_nl, 4l, 4l>::work_size(int, ggml_tensor const*, unsigned long&)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_iq4_nl, 4l, 4l>::~tensor_traits()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_iq4_nl, 4l, 4l>::~tensor_traits()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::tensor_traits<block_iq4_nl, 4l, 4l>::~tensor_traits()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::extra_buffer_type::supports_op(ggml_backend_device*, ggml_tensor const*)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::extra_buffer_type::get_tensor_traits(ggml_tensor const*)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::extra_buffer_type::~extra_buffer_type()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::extra_buffer_type::~extra_buffer_type()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml::cpu::aarch64::extra_buffer_type::~extra_buffer_type()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:void ggml::cpu::aarch64::gemm<block_q4_0, 4l, 4l>(int, float*, unsigned long, void const*, void const*, int, int)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:void ggml::cpu::aarch64::gemm<block_q4_0, 8l, 4l>(int, float*, unsigned long, void const*, void const*, int, int)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:void ggml::cpu::aarch64::gemm<block_q4_0, 8l, 8l>(int, float*, unsigned long, void const*, void const*, int, int)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:void ggml::cpu::aarch64::gemm<block_iq4_nl, 4l, 4l>(int, float*, unsigned long, void const*, void const*, int, int)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:void ggml::cpu::aarch64::gemv<block_q4_0, 4l, 4l>(int, float*, unsigned long, void const*, void const*, int, int)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:void ggml::cpu::aarch64::gemv<block_q4_0, 8l, 4l>(int, float*, unsigned long, void const*, void const*, int, int)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:void ggml::cpu::aarch64::gemv<block_q4_0, 8l, 8l>(int, float*, unsigned long, void const*, void const*, int, int)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:void ggml::cpu::aarch64::gemv<block_iq4_nl, 4l, 4l>(int, float*, unsigned long, void const*, void const*, int, int)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:int ggml::cpu::aarch64::repack<block_q4_0, 4l, 4l>(ggml_tensor*, void const*, unsigned long)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:int ggml::cpu::aarch64::repack<block_q4_0, 8l, 4l>(ggml_tensor*, void const*, unsigned long)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:int ggml::cpu::aarch64::repack<block_q4_0, 8l, 8l>(ggml_tensor*, void const*, unsigned long)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:int ggml::cpu::aarch64::repack<block_iq4_nl, 4l, 4l>(ggml_tensor*, void const*, unsigned long)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:cpuid_x86::cpuid_x86()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:cpuid_x86::cpuid_x86()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:void std::vector<ggml_backend_feature, std::allocator<ggml_backend_feature> >::_M_realloc_append<ggml_backend_feature>(ggml_backend_feature&&)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:std::vector<ggml_backend_feature, std::allocator<ggml_backend_feature> >::~vector()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:std::vector<ggml_backend_feature, std::allocator<ggml_backend_feature> >::~vector()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:std::vector<ggml_backend_buffer_type*, std::allocator<ggml_backend_buffer_type*> >::~vector()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:std::vector<ggml_backend_buffer_type*, std::allocator<ggml_backend_buffer_type*> >::~vector()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:void std::vector<std::array<int, 4ul>, std::allocator<std::array<int, 4ul> > >::_M_realloc_append<std::array<int, 4ul> const&>(std::array<int, 4ul> const&)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >::_M_dispose()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >::_M_replace_cold(char*, unsigned long, char const*, unsigned long, unsigned long)
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >::~basic_string()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >::~basic_string()
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_backend_cpu_init
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_backend_cpu_reg
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_backend_cpu_set_abort_callback
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_backend_cpu_set_n_threads
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_backend_cpu_set_threadpool
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_backend_init
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_backend_is_cpu
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_backend_score
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_barrier
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_extra_compute_forward
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_extra_work_size
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_get_sve_cnt
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_amx_int8
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_arm_fma
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_avx
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_avx2
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_avx512
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_avx512_bf16
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_avx512_vbmi
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_avx512_vnni
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_avx_vnni
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_dotprod
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_f16c
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_fma
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_fp16_va
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_llamafile
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_matmul_int8
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_neon
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_riscv_v
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_sme
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_sse3
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_ssse3
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_sve
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_vsx
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_vxe
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_has_wasm_simd
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_cpu_init
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_get_f32_1d
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_get_f32_nd
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_get_i32_1d
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_get_i32_nd
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_get_type_traits_cpu
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_graph_compute
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_graph_compute_with_ctx
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_graph_plan
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_is_numa
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_new_f32
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_new_i32
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_numa_init
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_set_f32
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_set_f32_1d
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_set_f32_nd
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_set_i32
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_set_i32_1d
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_set_i32_nd
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_threadpool_free
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_threadpool_new
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_threadpool_pause
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_threadpool_resume
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_iq1_m_q8_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_iq1_s_q8_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_iq2_s_q8_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_iq2_xs_q8_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_iq2_xxs_q8_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_iq3_s_q8_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_iq3_xxs_q8_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_iq4_nl_q8_0
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_iq4_xs_q8_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_q2_K_q8_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_q3_K_q8_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_q4_0_q8_0
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_q4_1_q8_1
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_q4_K_q8_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_q5_0_q8_0
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_q5_1_q8_1
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_q5_K_q8_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_q6_K_q8_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_q8_0_q8_0
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_tq1_0_q8_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:ggml_vec_dot_tq2_0_q8_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:llamafile_sgemm
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:quantize_row_iq4_nl
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:quantize_row_iq4_xs
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:quantize_row_q2_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:quantize_row_q3_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:quantize_row_q4_0
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:quantize_row_q4_1
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:quantize_row_q4_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:quantize_row_q5_0
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:quantize_row_q5_1
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:quantize_row_q5_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:quantize_row_q6_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:quantize_row_q8_0
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:quantize_row_q8_1
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:quantize_row_q8_K
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:quantize_row_tq1_0
uprobe:/usr/lib/ollama/libggml-cpu-alderlake.so:quantize_row_tq2_0


















ollama /usr/lib/ollama/libggml-cpu-haswell.so
ollama /usr/lib/ollama/libggml-cpu-icelake.so
ollama /usr/lib/ollama/libggml-cpu-sandybridge.so
ollama /usr/lib/ollama/libggml-cpu-skylakex.so



❯ paru -Ql ollama-cuda

ollama-cuda /usr/lib/ollama/cuda_v12/libggml-cuda.so
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 16, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 16, 1, 1, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 16, 1, 2, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 16, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 16, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 16, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 1, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 2, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 2, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 32, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 32, 1, 1, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 32, 1, 2, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 32, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 32, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 4, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 4, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 4, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 64, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 8, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 8, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 8, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<112, 8, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 16, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 16, 1, 1, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 16, 1, 2, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 16, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 16, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 16, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 1, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 1, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 2, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 2, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 2, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 32, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 32, 1, 1, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 32, 1, 2, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 32, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 32, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 4, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 4, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 4, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 4, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 64, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 8, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 8, 1, 1, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 8, 1, 2, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 8, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 8, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 8, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<128, 8, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 16, 1, 0, 32>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 16, 1, 1, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 16, 1, 2, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 16, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 16, 2, 0, 32>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 16, 4, 0, 32>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 1, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 1, 8, 0, 32>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 2, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 2, 4, 0, 32>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 2, 8, 0, 32>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 32, 1, 0, 32>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 32, 1, 1, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 32, 1, 2, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 32, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 32, 2, 0, 32>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 4, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 4, 2, 0, 32>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 4, 4, 0, 32>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 4, 8, 0, 32>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 64, 1, 0, 32>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 8, 1, 0, 32>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 8, 1, 1, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 8, 1, 2, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 8, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 8, 2, 0, 32>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 8, 4, 0, 32>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<256, 8, 8, 0, 32>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 16, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 16, 1, 1, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 16, 1, 2, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 16, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 16, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 16, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 1, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 1, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 2, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 2, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 2, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 32, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 32, 1, 1, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 32, 1, 2, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 32, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 32, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 4, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 4, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 4, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 4, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 64, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 8, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 8, 1, 1, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 8, 1, 2, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 8, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 8, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 8, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<64, 8, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 16, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 16, 1, 1, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 16, 1, 2, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 16, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 16, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 16, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 1, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 2, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 2, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 32, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 32, 1, 1, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 32, 1, 2, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 32, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 32, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 4, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 4, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 4, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 64, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 8, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 8, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 8, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<80, 8, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 16, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 16, 1, 1, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 16, 1, 2, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 16, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 16, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 16, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 1, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 2, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 2, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 32, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 32, 1, 1, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 32, 1, 2, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 32, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 32, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 4, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 4, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 4, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 64, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 8, 1, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 8, 1, 1, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 8, 1, 2, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 8, 1, 4, -1>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 8, 2, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 8, 4, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void launch_fattn<96, 8, 8, 0, 64>(ggml_backend_cuda_context&, ggml_tensor*, void (*)(char const*, char const*, char const*, char const*, float*, float2*, float, float, float, float, unsigned int, float, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int), int, unsigned long, bool, bool)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:sum_f32_cuda(ggml_cuda_pool&, float const*, float*, long, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_cpy(ggml_backend_cuda_context&, ggml_tensor const*, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_dup(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_info()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void mul_mat_q_case<(ggml_type)10>(ggml_backend_cuda_context&, mmq_args const&, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void mul_mat_q_case<(ggml_type)11>(ggml_backend_cuda_context&, mmq_args const&, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void mul_mat_q_case<(ggml_type)12>(ggml_backend_cuda_context&, mmq_args const&, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void mul_mat_q_case<(ggml_type)13>(ggml_backend_cuda_context&, mmq_args const&, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void mul_mat_q_case<(ggml_type)14>(ggml_backend_cuda_context&, mmq_args const&, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void mul_mat_q_case<(ggml_type)16>(ggml_backend_cuda_context&, mmq_args const&, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void mul_mat_q_case<(ggml_type)17>(ggml_backend_cuda_context&, mmq_args const&, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void mul_mat_q_case<(ggml_type)18>(ggml_backend_cuda_context&, mmq_args const&, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void mul_mat_q_case<(ggml_type)19>(ggml_backend_cuda_context&, mmq_args const&, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void mul_mat_q_case<(ggml_type)20>(ggml_backend_cuda_context&, mmq_args const&, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void mul_mat_q_case<(ggml_type)21>(ggml_backend_cuda_context&, mmq_args const&, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void mul_mat_q_case<(ggml_type)22>(ggml_backend_cuda_context&, mmq_args const&, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void mul_mat_q_case<(ggml_type)23>(ggml_backend_cuda_context&, mmq_args const&, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void mul_mat_q_case<(ggml_type)2>(ggml_backend_cuda_context&, mmq_args const&, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void mul_mat_q_case<(ggml_type)3>(ggml_backend_cuda_context&, mmq_args const&, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void mul_mat_q_case<(ggml_type)6>(ggml_backend_cuda_context&, mmq_args const&, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void mul_mat_q_case<(ggml_type)7>(ggml_backend_cuda_context&, mmq_args const&, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void mul_mat_q_case<(ggml_type)8>(ggml_backend_cuda_context&, mmq_args const&, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_error(char const*, char const*, char const*, int, char const*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_argmax(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_cpy_fn(ggml_tensor const*, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_acc(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_add(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_cos(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_div(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_exp(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_mul(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_neg(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_pad(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_sin(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_sqr(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_sub(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_sum(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_gelu(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_norm(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_relu(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_rope(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_silu(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_sqrt(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_step(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_tanh(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:sum_rows_f32_cuda(float const*, float*, int, int, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_clamp(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_scale(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_unpad(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_out_prod(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_arange(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_concat(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_im2col(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_pool2d(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_repeat(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_get_device()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_argsort(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_sigmoid(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_upscale(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_set_device(int)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_count_equal(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_mul_mat_vec(ggml_backend_cuda_context&, ggml_tensor const*, ggml_tensor const*, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_get_rows(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_rms_norm(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_soft_max(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_sum_rows(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_get_to_fp16_cuda(ggml_type)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_get_to_fp32_cuda(ggml_type)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_hardswish(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_mul_mat_q(ggml_backend_cuda_context&, ggml_tensor const*, ggml_tensor const*, ggml_tensor*, char const*, float const*, char const*, float*, long, long, long, long, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_rope_back(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_op_rope_impl<false>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_op_rope_impl<true>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_rwkv_wkv6(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_silu_back(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:quantize_mmq_q8_1_cuda(float const*, void*, long, long, long, long, ggml_type, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:quantize_row_q8_1_cuda(float const*, void*, long, long, long, long, ggml_type, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_gelu_quick(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_group_norm(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_leaky_relu(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_flash_attn_ext(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_hardsigmoid(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_mul_mat_vec(ggml_backend_cuda_context&, ggml_tensor const*, ggml_tensor const*, ggml_tensor*, char const*, float const*, char const*, float*, long, long, long, long, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_repeat_back(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_opt_step_adamw(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_should_use_mmq(ggml_type, int, long)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_diag_mask_inf(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_get_rows_back(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_mul_mat_vec_q(ggml_backend_cuda_context&, ggml_tensor const*, ggml_tensor const*, ggml_tensor*, char const*, float const*, char const*, float*, long, long, long, long, CUstream_st*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_rms_norm_back(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_soft_max_back(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_cross_entropy_loss(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_conv_transpose_1d(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_gated_linear_attn(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_highest_compiled_arch(int)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_op_timestep_embedding(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_cross_entropy_loss_back(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_flash_attn_ext_tile_f16(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_flash_attn_ext_tile_f32(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_flash_attn_ext_wmma_f16(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:int ggml_cuda_highest_compiled_arch_impl<int, int, int, int>(int, int, int, int, int, int, int)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:int ggml_cuda_highest_compiled_arch_impl<int, int, int, int, int, int, int, int>(int, int, int, int, int, int, int, int, int, int, int)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:int ggml_cuda_highest_compiled_arch_impl<int, int, int, int, int, int, int, int, int, int, int, int>(int, int, int, int, int, int, int, int, int, int, int, int, int, int, int)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<112, 16, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<112, 16, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<112, 16, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<112, 1, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<112, 2, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<112, 2, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<112, 32, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<112, 32, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<112, 4, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<112, 4, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<112, 4, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<112, 64, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<112, 8, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<112, 8, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<112, 8, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<112, 8, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<128, 16, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<128, 16, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<128, 16, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<128, 1, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<128, 2, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<128, 2, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<128, 32, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<128, 32, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<128, 4, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<128, 4, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<128, 4, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<128, 64, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<128, 8, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<128, 8, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<128, 8, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<128, 8, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<256, 16, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<256, 16, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<256, 16, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<256, 1, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<256, 2, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<256, 2, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<256, 32, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<256, 32, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<256, 4, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<256, 4, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<256, 4, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<256, 64, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<256, 8, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<256, 8, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<256, 8, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<256, 8, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<64, 16, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<64, 16, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<64, 16, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<64, 1, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<64, 2, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<64, 2, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<64, 32, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<64, 32, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<64, 4, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<64, 4, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<64, 4, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<64, 64, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<64, 8, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<64, 8, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<64, 8, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<64, 8, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<80, 16, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<80, 16, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<80, 16, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<80, 1, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<80, 2, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<80, 2, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<80, 32, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<80, 32, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<80, 4, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<80, 4, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<80, 4, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<80, 64, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<80, 8, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<80, 8, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<80, 8, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<80, 8, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<96, 16, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<96, 16, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<96, 16, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<96, 1, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<96, 2, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<96, 2, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<96, 32, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<96, 32, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<96, 4, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<96, 4, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<96, 4, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<96, 64, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<96, 8, 1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<96, 8, 2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<96, 8, 4>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_mma_f16_case<96, 8, 8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_vec_f16_case<128, (ggml_type)1, (ggml_type)1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_vec_f16_case<128, (ggml_type)2, (ggml_type)2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_vec_f16_case<128, (ggml_type)8, (ggml_type)8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_vec_f16_case<256, (ggml_type)1, (ggml_type)1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_vec_f16_case<64, (ggml_type)1, (ggml_type)1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_vec_f32_case<128, (ggml_type)1, (ggml_type)1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_vec_f32_case<128, (ggml_type)2, (ggml_type)2>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_vec_f32_case<128, (ggml_type)8, (ggml_type)8>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_vec_f32_case<256, (ggml_type)1, (ggml_type)1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_vec_f32_case<64, (ggml_type)1, (ggml_type)1>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<112, 16, __half>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<112, 16, float>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<112, 32, __half>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<112, 32, float>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<128, 16, __half>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<128, 16, float>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<128, 32, __half>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<128, 32, float>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<128, 8, __half>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<256, 16, __half>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<256, 16, float>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<256, 32, __half>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<256, 8, __half>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<64, 16, __half>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<64, 16, float>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<64, 32, __half>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<64, 32, float>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<64, 8, __half>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<80, 16, __half>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<80, 16, float>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<80, 32, __half>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<80, 32, float>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<96, 16, __half>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<96, 16, float>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<96, 32, __half>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<96, 32, float>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void ggml_cuda_flash_attn_ext_wmma_f16_case<96, 8, __half>(ggml_backend_cuda_context&, ggml_tensor*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_graph::~ggml_cuda_graph()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_graph::~ggml_cuda_graph()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_pool_leg::free(void*, unsigned long)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_pool_leg::alloc(unsigned long, unsigned long*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_pool_leg::~ggml_cuda_pool_leg()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_pool_leg::~ggml_cuda_pool_leg()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_pool_leg::~ggml_cuda_pool_leg()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_pool_vmm::free(void*, unsigned long)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_pool_vmm::alloc(unsigned long, unsigned long*)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_pool_vmm::~ggml_cuda_pool_vmm()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_pool_vmm::~ggml_cuda_pool_vmm()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_cuda_pool_vmm::~ggml_cuda_pool_vmm()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_backend_cuda_context::cublas_handle(int)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_backend_cuda_context::new_pool_for_device(int)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_backend_cuda_context::pool(int)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:cuda::__4::cuda_error::~cuda_error()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:cuda::__4::cuda_error::~cuda_error()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:cuda::__4::cuda_error::~cuda_error()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:nvtx3::v1::scoped_range_in<cub::CUB_200700_500_520_530_600_610_620_700_720_750_800_860_870_890_900_NS::detail::NVTXCCCLDomain>::~scoped_range_in()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:nvtx3::v1::scoped_range_in<cub::CUB_200700_500_520_530_600_610_620_700_720_750_800_860_870_890_900_NS::detail::NVTXCCCLDomain>::~scoped_range_in()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:std::map<std::pair<int, std::array<float, 16ul> >, ggml_backend_buffer_type, std::less<std::pair<int, std::array<float, 16ul> > >, std::allocator<std::pair<std::pair<int, std::array<float, 16ul> > const, ggml_backend_buffer_type> > >::~map()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:std::map<std::pair<int, std::array<float, 16ul> >, ggml_backend_buffer_type, std::less<std::pair<int, std::array<float, 16ul> > >, std::allocator<std::pair<std::pair<int, std::array<float, 16ul> > const, ggml_backend_buffer_type> > >::~map()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:std::vector<cudaKernelNodeParams, std::allocator<cudaKernelNodeParams> >::_M_default_append(unsigned long)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void std::vector<ggml_backend_feature, std::allocator<ggml_backend_feature> >::_M_realloc_insert<ggml_backend_feature>(__gnu_cxx::__normal_iterator<ggml_backend_feature*, std::vector<ggml_backend_feature, std::allocator<ggml_backend_feature> > >, ggml_backend_feature&&)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:std::vector<ggml_backend_feature, std::allocator<ggml_backend_feature> >::~vector()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:std::vector<ggml_backend_feature, std::allocator<ggml_backend_feature> >::~vector()
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:std::vector<ggml_graph_node_properties, std::allocator<ggml_graph_node_properties> >::_M_default_append(unsigned long)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:std::vector<CUgraphNode_st*, std::allocator<CUgraphNode_st*> >::_M_default_append(unsigned long)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void std::vector<ggml_backend_device*, std::allocator<ggml_backend_device*> >::_M_realloc_insert<ggml_backend_device* const&>(__gnu_cxx::__normal_iterator<ggml_backend_device**, std::vector<ggml_backend_device*, std::allocator<ggml_backend_device*> > >, ggml_backend_device* const&)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void std::vector<ggml_tensor_extra_gpu*, std::allocator<ggml_tensor_extra_gpu*> >::_M_realloc_insert<ggml_tensor_extra_gpu* const&>(__gnu_cxx::__normal_iterator<ggml_tensor_extra_gpu**, std::vector<ggml_tensor_extra_gpu*, std::allocator<ggml_tensor_extra_gpu*> > >, ggml_tensor_extra_gpu* const&)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void std::vector<char**, std::allocator<char**> >::_M_realloc_insert<char**>(__gnu_cxx::__normal_iterator<char***, std::vector<char**, std::allocator<char**> > >, char**&&)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:void std::vector<void*, std::allocator<void*> >::_M_realloc_insert<void* const&>(__gnu_cxx::__normal_iterator<void**, std::vector<void*, std::allocator<void*> > >, void* const&)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:std::__cxx11::to_string(int)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:std::_Rb_tree_iterator<std::pair<std::pair<int, std::array<float, 16ul> > const, ggml_backend_buffer_type> > std::_Rb_tree<std::pair<int, std::array<float, 16ul> >, std::pair<std::pair<int, std::array<float, 16ul> > const, ggml_backend_buffer_type>, std::_Select1st<std::pair<std::pair<int, std::array<float, 16ul> > const, ggml_backend_buffer_type> >, std::less<std::pair<int, std::array<float, 16ul> > >, std::allocator<std::pair<std::pair<int, std::array<float, 16ul> > const, ggml_backend_buffer_type> > >::_M_emplace_hint_unique<std::pair<int, std::array<float, 16ul> >, ggml_backend_buffer_type&>(std::_Rb_tree_const_iterator<std::pair<std::pair<int, std::array<float, 16ul> > const, ggml_backend_buffer_type> >, std::pair<int, std::array<float, 16ul> >&&, ggml_backend_buffer_type&)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:std::_Rb_tree<std::pair<int, std::array<float, 16ul> >, std::pair<std::pair<int, std::array<float, 16ul> > const, ggml_backend_buffer_type>, std::_Select1st<std::pair<std::pair<int, std::array<float, 16ul> > const, ggml_backend_buffer_type> >, std::less<std::pair<int, std::array<float, 16ul> > >, std::allocator<std::pair<std::pair<int, std::array<float, 16ul> > const, ggml_backend_buffer_type> > >::_M_get_insert_unique_pos(std::pair<int, std::array<float, 16ul> > const&)
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_backend_cuda_buffer_type
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_backend_cuda_get_device_count
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_backend_cuda_get_device_description
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_backend_cuda_get_device_memory
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_backend_cuda_host_buffer_type
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_backend_cuda_init
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_backend_cuda_reg
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_backend_cuda_register_host_buffer
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_backend_cuda_split_buffer_type
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_backend_cuda_unregister_host_buffer
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_backend_init
uprobe:/usr/lib/ollama/cuda_v12/libggml-cuda.so:ggml_backend_is_cuda
