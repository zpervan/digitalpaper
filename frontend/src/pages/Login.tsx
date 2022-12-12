import { useState } from 'react';
import Layout from '../components/Layout';
import axios from 'axios';

const Login = () => {
  const [data, setData] = useState({
    username: '',
    mail: '',
    password: '',
  });

  const getUser = (e: React.FormEvent, url: string) => {
    e.preventDefault();

    axios.get(url).then((res) => console.log(res));

    setData({
      username: '',
      mail: '',
      password: '',
    });
  };

  return (
    <Layout>
      <main className='flex justify-center'>
        <section className='w-1/2 flex flex-col justify-center'>
          <h2 className='text-white text-3xl mb-6 font-extrabold'>Sign Up</h2>
          <form
            className='flex flex-col mt-4 text-main-white'
            onSubmit={(e) => getUser(e, 'http://localhost:3500/users')}
          >
            <input
              type='text'
              className='bg-zinc-800 rounded-md px-4 py-4 my-3 h-12'
              placeholder='Username'
              onChange={(e) => {
                setData({ ...data, username: e.target.value });
              }}
              value={data.username}
            />
            <input
              type='email'
              placeholder='Email'
              className='bg-zinc-800 rounded-md px-4 py-4 my-3 h-12'
              onChange={(e) => {
                setData({ ...data, mail: e.target.value });
              }}
              value={data.mail}
            />
            <input
              type='password'
              placeholder='Password'
              className='bg-zinc-800 rounded-md px-4 py-4 my-3 h-12'
              onChange={(e) => {
                setData({ ...data, password: e.target.value });
              }}
              value={data.password}
            />
            <button
              type='submit'
              className='bg-main-red py-3 mt-4 rounded-sm text-main-white text-base font-bold uppercase'
            >
              Sign up
            </button>
          </form>
        </section>
      </main>
    </Layout>
  );
};

export default Login;
