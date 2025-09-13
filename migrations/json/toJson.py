import pandas as pd

df = pd.read_csv('migrations/json/Pendataan Jemaat ENST.csv')

print(df["Tanggal Lahir"])
# print(df["Spiritual Journey"].unique())
# ['Timestamp', 'Nama Lengkap', 'Nama Lain', 'Gender', 'Tempat Lahir',
#        'Tanggal Lahir', 'Fase Hidup', 'Status Perkawinan', 'Alamat',
#        'Nomor Telepon (WA)', 'Pemimpin Life Group', 'Email', 'Dimuridkan Oleh',
#        'Spiritual Journey', 'Tanggal Pernikahan', 'Nama Pasangan']

data = {
    'nama': df['Nama Lengkap'].astype('string'),
    'nama_lain': df['Nama Lain'].astype('string'),
    'gender': df['Gender'].replace({'Laki-laki': 'L', 'Perempuan': 'P'}).astype('string'),
    'tempat_lahir': df['Tempat Lahir'].astype('string'),
    'tanggal_lahir': pd.to_datetime(df['Tanggal Lahir']),
    'fase_hidup': df['Fase Hidup'].astype('string'),
    'status_perkawinan': df['Status Perkawinan'].astype('string'),
    'nama_pasangan': df['Nama Pasangan'].astype('string'),
    'tanggal_perkawinan': pd.to_datetime(df['Tanggal Pernikahan'], errors='coerce'),
    'alamat': df['Alamat'].astype('string'),
    'nomor_telepon': df['Nomor Telepon (WA)'].astype('string'),
    'email': df['Email'].astype('string'),
    'church_id': "123e4567-e89b-12d3-a456-426614174002",
    'kabupaten_id': 3578
}


new_df = pd.DataFrame(data)
new_df['tanggal_lahir'] = new_df['tanggal_lahir'].dt.strftime('%Y-%m-%dT%H:%M:%SZ')
new_df['tanggal_perkawinan'] = new_df['tanggal_perkawinan'].dt.strftime('%Y-%m-%dT%H:%M:%SZ')
print(new_df.head())

new_df.to_json('migrations/json/person.json', orient='records', lines=True)