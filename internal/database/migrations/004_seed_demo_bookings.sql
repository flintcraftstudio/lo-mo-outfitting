-- Demo booking requests for UI development
INSERT INTO booking_requests (
    ip_address, trip_type, preferred_date, alternate_date, angler_count, youth_count, heroes,
    experience, lodging, lodging_other, client_notes, referred_by,
    client_name, client_email, client_phone, status, guide_id, payment_method, mat_notes, status_updated_at, emailed_at, created_at
) VALUES
-- New submissions
('192.168.1.10', 'full_day_single', '2026-05-15', '2026-05-16', '2', '0', 0,
 'comfortable', 'craig', '', 'We fish together every year. Looking forward to the Missouri!', 'Google search',
 'Jake Morrison', 'jake.morrison@gmail.com', '406-555-1234', 'new', NULL, NULL, '', NULL, '2026-04-07 14:30:00', '2026-04-07 14:30:00'),

('10.0.0.5', 'heroes', '2026-05-20', '', '1', '0', 1,
 'some', 'helena', '', 'Army veteran, stationed at Malmstrom. First time on the Missouri.', 'Buddy at the base',
 'Sgt. Marcus Bell', 'mbell_mt@yahoo.com', '406-555-9821', 'new', NULL, NULL, '', NULL, '2026-04-08 09:15:00', '2026-04-08 09:15:00'),

('172.16.0.22', 'half_day_single', '2026-04-25', '', '3', '1', 0,
 'never', 'not_sure', '', 'Family trip — kids are 14 and 10. The 14-year-old is super excited.', 'TripAdvisor',
 'Rachel Nguyen', 'rachel.nguyen@outlook.com', '303-555-4567', 'new', NULL, NULL, '', NULL, '2026-04-08 11:00:00', '2026-04-08 11:00:00'),

-- Contacted
('192.168.1.44', 'full_day_single', '2026-05-10', '2026-05-11', '2', '0', 0,
 'advanced', 'craig', '', 'Dry fly purists. Happy to wade if the flows are right.', 'Montana Fly Company shop',
 'Tom Bradshaw', 'tbradshaw@troutmail.com', '208-555-3399', 'contacted', NULL, NULL, 'Called 4/6 — wants to confirm lodging first. Follow up Friday.', '2026-04-06 16:00:00', '2026-04-04 10:20:00', '2026-04-04 10:20:00'),

('10.10.0.8', 'multiple_boats', '2026-06-14', '2026-06-15', '6', '0', 0,
 'comfortable', 'wolf_creek', '', 'Bachelor party — 6 guys, 3 boats. We want the full experience.', 'Instagram',
 'Derek Haines', 'dhaines@proton.me', '503-555-7712', 'contacted', NULL, NULL, 'Big group. Need 3 guides. Check availability with Colter and Dylan.', '2026-04-05 11:00:00', '2026-04-03 08:45:00', '2026-04-03 08:45:00'),

-- Deposit sent
('192.168.2.100', 'early_season', '2026-04-18', '', '2', '0', 0,
 'comfortable', 'great_falls', '', '', 'Return client — fished with Ria last year',
 'Bill Westbrook', 'bill.westbrook@me.com', '406-555-6001', 'deposit_sent', 1, NULL, 'Repeat client. Loves Ria. Venmo sent, waiting to clear.', '2026-04-02 09:30:00', '2026-03-28 15:00:00', '2026-03-28 15:00:00'),

-- Confirmed
('10.0.1.15', 'full_day_single', '2026-04-12', '', '2', '0', 0,
 'advanced', 'craig', '', 'Streamer day if possible. Big bug water.', 'Headhunters fly shop',
 'Kyle Fenton', 'kfenton@sportfish.net', '307-555-2288', 'confirmed', 2, 'venmo', 'Confirmed. Colter guiding. Streamer setup ready.', '2026-04-01 14:00:00', '2026-03-25 12:00:00', '2026-03-25 12:00:00'),

('172.16.5.3', 'full_day_single', '2026-04-19', '2026-04-20', '1', '0', 0,
 'some', 'helena', '', 'Solo angler, visiting from Seattle. Open to whatever is fishing well.', 'Orvis website referral',
 'Megan Park', 'mpark.seattle@gmail.com', '206-555-4410', 'confirmed', 3, 'stripe', '', '2026-04-03 10:00:00', '2026-03-30 17:30:00', '2026-03-30 17:30:00'),

('192.168.0.77', 'heroes', '2026-04-26', '', '2', '0', 1,
 'comfortable', 'craig', '', 'Retired fire captain and my son. He just got back from wildfire season.', 'Lo Mo website',
 'Capt. Dan Reeves', 'dreeves.fire@gmail.com', '406-555-8833', 'confirmed', 4, NULL, 'Heroes rate approved. Dylan guiding. No deposit yet — follow up.', '2026-04-04 08:00:00', '2026-03-31 11:15:00', '2026-03-31 11:15:00'),

('10.0.0.99', 'full_day_single', '2026-05-03', '', '2', '0', 0,
 'comfortable', 'wolf_creek', '', 'Anniversary trip with my wife. She outfishes me every time.', 'Friend recommendation',
 'Steve Yamamoto', 'stevey@gmail.com', '415-555-1190', 'confirmed', 5, 'cash', 'Cash deposit received at the shop.', '2026-04-05 13:00:00', '2026-04-01 09:00:00', '2026-04-01 09:00:00'),

-- Complete
('192.168.1.200', 'full_day_single', '2026-03-22', '', '2', '0', 0,
 'advanced', 'craig', '', 'Spring BWO hatch — hoping for good timing.', 'Repeat client',
 'Frank Deluca', 'fdeluca@fishon.com', '406-555-3344', 'complete', 1, 'venmo', 'Great day on the water. Landed a 22" brown on a BWO emerger.', '2026-03-22 18:00:00', '2026-03-10 14:00:00', '2026-03-10 14:00:00'),

-- Cancelled
('10.10.10.10', 'half_day_single', '2026-04-15', '', '2', '1', 0,
 'never', 'helena', '', 'First time fly fishing!', 'Yelp',
 'Amy Chen', 'amychen88@gmail.com', '406-555-7700', 'cancelled', NULL, NULL, 'Cancelled — family emergency. Wants to rebook later this summer.', '2026-04-06 10:00:00', '2026-04-02 16:30:00', '2026-04-02 16:30:00');

-- Events for the demo bookings
-- Jake Morrison (id 1) - new
INSERT INTO booking_events (booking_request_id, event_type, detail, created_at) VALUES
(1, 'submitted', 'Booking request submitted via website', '2026-04-07 14:30:00'),
(1, 'email_sent', 'Notification email sent', '2026-04-07 14:30:01');

-- Sgt. Marcus Bell (id 2) - new
INSERT INTO booking_events (booking_request_id, event_type, detail, created_at) VALUES
(2, 'submitted', 'Booking request submitted via website', '2026-04-08 09:15:00'),
(2, 'email_sent', 'Notification email sent', '2026-04-08 09:15:01');

-- Rachel Nguyen (id 3) - new
INSERT INTO booking_events (booking_request_id, event_type, detail, created_at) VALUES
(3, 'submitted', 'Booking request submitted via website', '2026-04-08 11:00:00'),
(3, 'email_sent', 'Notification email sent', '2026-04-08 11:00:01');

-- Tom Bradshaw (id 4) - contacted
INSERT INTO booking_events (booking_request_id, event_type, detail, created_at) VALUES
(4, 'submitted', 'Booking request submitted via website', '2026-04-04 10:20:00'),
(4, 'email_sent', 'Notification email sent', '2026-04-04 10:20:01'),
(4, 'status_changed', 'Status changed: new → contacted', '2026-04-06 16:00:00'),
(4, 'note_added', 'Note: Called 4/6 — wants to confirm lodging first. Follow up Friday.', '2026-04-06 16:05:00');

-- Derek Haines (id 5) - contacted
INSERT INTO booking_events (booking_request_id, event_type, detail, created_at) VALUES
(5, 'submitted', 'Booking request submitted via website', '2026-04-03 08:45:00'),
(5, 'email_sent', 'Notification email sent', '2026-04-03 08:45:01'),
(5, 'status_changed', 'Status changed: new → contacted', '2026-04-05 11:00:00'),
(5, 'note_added', 'Note: Big group. Need 3 guides. Check availability with Colter and Dylan.', '2026-04-05 11:05:00');

-- Bill Westbrook (id 6) - deposit_sent
INSERT INTO booking_events (booking_request_id, event_type, detail, created_at) VALUES
(6, 'submitted', 'Booking request submitted via website', '2026-03-28 15:00:00'),
(6, 'email_sent', 'Notification email sent', '2026-03-28 15:00:01'),
(6, 'status_changed', 'Status changed: new → contacted', '2026-03-29 10:00:00'),
(6, 'guide_assigned', 'Guide assigned: Ria French', '2026-03-29 10:05:00'),
(6, 'status_changed', 'Status changed: contacted → deposit_sent', '2026-04-02 09:30:00'),
(6, 'note_added', 'Note: Repeat client. Loves Ria. Venmo sent, waiting to clear.', '2026-04-02 09:35:00');

-- Kyle Fenton (id 7) - confirmed
INSERT INTO booking_events (booking_request_id, event_type, detail, created_at) VALUES
(7, 'submitted', 'Booking request submitted via website', '2026-03-25 12:00:00'),
(7, 'email_sent', 'Notification email sent', '2026-03-25 12:00:01'),
(7, 'status_changed', 'Status changed: new → contacted', '2026-03-26 09:00:00'),
(7, 'guide_assigned', 'Guide assigned: Colter Day', '2026-03-26 09:05:00'),
(7, 'status_changed', 'Status changed: contacted → deposit_sent', '2026-03-28 14:00:00'),
(7, 'payment_method_set', 'Payment method set: venmo', '2026-03-28 14:05:00'),
(7, 'status_changed', 'Status changed: deposit_sent → confirmed', '2026-04-01 14:00:00');

-- Megan Park (id 8) - confirmed
INSERT INTO booking_events (booking_request_id, event_type, detail, created_at) VALUES
(8, 'submitted', 'Booking request submitted via website', '2026-03-30 17:30:00'),
(8, 'email_sent', 'Notification email sent', '2026-03-30 17:30:01'),
(8, 'status_changed', 'Status changed: new → contacted', '2026-03-31 08:00:00'),
(8, 'status_changed', 'Status changed: contacted → deposit_sent', '2026-04-01 10:00:00'),
(8, 'guide_assigned', 'Guide assigned: Andrew Osborn', '2026-04-01 10:05:00'),
(8, 'payment_method_set', 'Payment method set: stripe', '2026-04-01 10:10:00'),
(8, 'status_changed', 'Status changed: deposit_sent → confirmed', '2026-04-03 10:00:00');

-- Capt. Dan Reeves (id 9) - confirmed, no deposit
INSERT INTO booking_events (booking_request_id, event_type, detail, created_at) VALUES
(9, 'submitted', 'Booking request submitted via website', '2026-03-31 11:15:00'),
(9, 'email_sent', 'Notification email sent', '2026-03-31 11:15:01'),
(9, 'status_changed', 'Status changed: new → contacted', '2026-04-01 09:00:00'),
(9, 'guide_assigned', 'Guide assigned: Dylan Huseby', '2026-04-01 09:05:00'),
(9, 'status_changed', 'Status changed: contacted → confirmed', '2026-04-04 08:00:00'),
(9, 'note_added', 'Note: Heroes rate approved. Dylan guiding. No deposit yet — follow up.', '2026-04-04 08:05:00');

-- Steve Yamamoto (id 10) - confirmed
INSERT INTO booking_events (booking_request_id, event_type, detail, created_at) VALUES
(10, 'submitted', 'Booking request submitted via website', '2026-04-01 09:00:00'),
(10, 'email_sent', 'Notification email sent', '2026-04-01 09:00:01'),
(10, 'status_changed', 'Status changed: new → contacted', '2026-04-02 08:00:00'),
(10, 'guide_assigned', 'Guide assigned: Sam Botz', '2026-04-02 08:05:00'),
(10, 'status_changed', 'Status changed: contacted → deposit_sent', '2026-04-03 14:00:00'),
(10, 'payment_method_set', 'Payment method set: cash', '2026-04-03 14:05:00'),
(10, 'status_changed', 'Status changed: deposit_sent → confirmed', '2026-04-05 13:00:00'),
(10, 'note_added', 'Note: Cash deposit received at the shop.', '2026-04-05 13:05:00');

-- Frank Deluca (id 11) - complete
INSERT INTO booking_events (booking_request_id, event_type, detail, created_at) VALUES
(11, 'submitted', 'Booking request submitted via website', '2026-03-10 14:00:00'),
(11, 'email_sent', 'Notification email sent', '2026-03-10 14:00:01'),
(11, 'status_changed', 'Status changed: new → contacted', '2026-03-11 09:00:00'),
(11, 'guide_assigned', 'Guide assigned: Ria French', '2026-03-11 09:05:00'),
(11, 'status_changed', 'Status changed: contacted → confirmed', '2026-03-15 10:00:00'),
(11, 'payment_method_set', 'Payment method set: venmo', '2026-03-15 10:05:00'),
(11, 'status_changed', 'Status changed: confirmed → complete', '2026-03-22 18:00:00'),
(11, 'note_added', 'Note: Great day on the water. Landed a 22" brown on a BWO emerger.', '2026-03-22 18:05:00');

-- Amy Chen (id 12) - cancelled
INSERT INTO booking_events (booking_request_id, event_type, detail, created_at) VALUES
(12, 'submitted', 'Booking request submitted via website', '2026-04-02 16:30:00'),
(12, 'email_sent', 'Notification email sent', '2026-04-02 16:30:01'),
(12, 'status_changed', 'Status changed: new → contacted', '2026-04-03 10:00:00'),
(12, 'status_changed', 'Status changed: contacted → cancelled', '2026-04-06 10:00:00'),
(12, 'note_added', 'Note: Cancelled — family emergency. Wants to rebook later this summer.', '2026-04-06 10:05:00');
